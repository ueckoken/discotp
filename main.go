package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env/v6"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Config struct {
	DiscordToken    string   `env:"DISCORD_APP_TOKEN"`
	AllowChannelIDs []string `env:"ALLOWED_REPLY_CHANNEL_IDS"`
	DiscordGuildID  string   `env:"DISCORD_GUILD_ID"`
	IsDevelopment   bool     `env:"IS_DEVELOPMENT" envDefault:"false"`
	Tokens          TotpToks `env:"TOTP_TOKENS"`
}

type totpTok string
type service string
type TotpToks struct {
	m map[service]totpTok
}

func (t *TotpToks) UnmarshalText(text []byte) error {
	type kv struct {
		k service
		v totpTok
	}
	ress := lo.Map(strings.Split(string(text), `,`), func(item string, index int) kv {
		ss := lo.Map(strings.SplitN(item, `:`, 2), func(item string, index int) string { return strings.TrimSpace(item) })
		return kv{k: service(ss[0]), v: totpTok(ss[1])}
	})
	t.m = lo.Associate(ress, func(item kv) (service, totpTok) { return item.k, item.v })
	return nil
}
func (s service) String() string { return string(s) }

func main() {
	config := &Config{}
	if err := env.Parse(config, env.Options{RequiredIfNoDef: true}); err != nil {
		log.Print(err)
		log.Fatalln("parse config failed")
	}
	logger := func(isDevel bool) *zap.Logger {
		if isDevel {
			return zap.Must(zap.NewDevelopment())
		}
		return zap.Must(zap.NewProduction())
	}(config.IsDevelopment)
	logger.Info("logger init finish", zap.Stringer("loglevel", logger.Level()))
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

	totpHandler := NewTotpHandler(logger, *config)
	appCmd, err := totpHandler.CreateTotpApplicationCommand(discordSession, config.DiscordGuildID)
	if err != nil {
		logger.Fatal("message create failed", zap.Error(err), zap.String("apply command", appCmd.Name))
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
