package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pquerna/otp/totp"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type TotpHandler struct {
	logger *zap.Logger
	config Config
}

func NewApplicationCommand(s *discordgo.Session, guildID string) (*discordgo.ApplicationCommand, error) {
	return s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "2fa",
		Description: "get 2fa code",
		Type:        discordgo.ChatApplicationCommand,

		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "google",
				Description: "get google 2fa code",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	})
}

func (h *TotpHandler) HandleIntractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var totpClient TotpClient
	switch i.ApplicationCommandData().Options[0].Name {
	case "google":
		totpClient = &GoogleTotp{secret: string(h.config.Google)}
	default:
		h.logger.Warn("sub command not found", zap.String("command name", i.ApplicationCommandData().Name))
		return fmt.Errorf("not exist sub command")
	}
	h.logger.Debug("channel IDS", zap.Strings("channel IDs", h.config.AllowChannelIDs), zap.Any("channel ID", i.ChannelID))
	if lo.Count(h.config.AllowChannelIDs, i.ChannelID) == 0 {
		h.logger.Debug("message hook in not allowed channel", zap.String("channelID", i.ChannelID))
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "このチャンネルからの呼出は許可されていません",
			},
		})
	}
	code, err := totpClient.GenerateCode(time.Now())
	if err != nil {
		h.logger.Warn("generate code error", zap.Error(err))
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Internal Error",
			},
		})
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("2faコードは`%s`です\n", code),
		},
	})
}

type TotpClient interface {
	GenerateCode(t time.Time) (string, error)
}

type GoogleTotp struct {
	secret string
}

func (g *GoogleTotp) GenerateCode(t time.Time) (string, error) {
	return totp.GenerateCode(g.secret, t)
}

func IntractionCreateHandlerRouter(logger *zap.Logger, config Config, cmdName string) (func(s *discordgo.Session, i *discordgo.InteractionCreate) error, error) {
	switch cmdName {
	case "2fa":
		handler := &TotpHandler{logger: logger, config: config}
		return handler.HandleIntractionCreate, nil
	default:
		return nil, fmt.Errorf("no route to handler")
	}
}
