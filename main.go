package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type Config struct {
	DiscordToken    string   `env:"DISCORD_APP_TOKEN"`
	AllowChannelIDs []string `env:"ALLOWED_REPLY_CHANNEL_IDS"`
	DiscordGuildID  string   `env:"DISCORD_GUILD_ID"`
	IsDevelopment   bool     `env:"IS_DEVELOPMENT" envDefault:"false"`
	Google          totpTok  `env:"GOOGLE_TOTP_TOKEN"`
}

type totpTok string

func main() {
	config := &Config{}
	if err := env.Parse(config, env.Options{RequiredIfNoDef: true}); err != nil {
		log.Print(err)
		log.Fatalln("parse config failed")
	}
	logger, err := func(isDevel bool) (*zap.Logger, error) {
		if isDevel {
			return zap.NewDevelopment()
		}
		return zap.NewProduction()
	}(config.IsDevelopment)
	if err != nil {
		log.Fatalln("logger init failed")
	}
	logger.Info("logger init finish")
	discordSession, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		logger.Fatal("discord go client init failed", zap.Error(err))
		return
	}

	if err := discordSession.Open(); err != nil {
		logger.Fatal("discord session open error", zap.Error(err))
		return
	}
	logger.Info("session opened")
	defer discordSession.Close()

	appCmd, err := NewApplicationCommand(discordSession, config.DiscordGuildID)
	if err != nil {
		logger.Fatal("message create failed", zap.Error(err))
		return
	}
	_, err = discordSession.ApplicationCommandCreate(discordSession.State.User.ID, config.DiscordGuildID, appCmd)
	if err != nil {
		logger.Error("application command create failed", zap.Error(err), zap.Any("command", appCmd))
	}
	discordSession.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handler, err := IntractionCreateHandlerRouter(logger, *config, i.ApplicationCommandData().Name)
		if err != nil {
			logger.Warn("interaction command not found error", zap.Error(err))
		}
		if err := handler(s, i); err != nil {
			logger.Warn("error in handler", zap.Error(err))
		}
	})
	logger.Info("add handler finish")
	logger.Info("start listening")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	logger.Info("stop signal handle")
}
