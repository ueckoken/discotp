package main

import (
	"fmt"
	"strings"
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

func NewTotpHandler(logger *zap.Logger, config Config) *TotpHandler {
	return &TotpHandler{
		logger: logger,
		config: config,
	}
}

func (h *TotpHandler) CreateTotpApplicationCommand(s *discordgo.Session, guildID string) (*discordgo.ApplicationCommand, error) {
	opts := lo.Map(lo.Keys(h.config.Tokens.m), func(item service, index int) *discordgo.ApplicationCommandOption {
		return &discordgo.ApplicationCommandOption{
			Name:        string(item),
			Description: fmt.Sprintf("%sの2faコードを取得します\n", item),
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		}
	})
	return s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "2fa",
		Description: "2faコードを取得することができます",
		Type:        discordgo.ChatApplicationCommand,
		Options:     opts,
	})
}

func (h *TotpHandler) HandleIntractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	h.logger.Debug("channel IDs", zap.Strings("channel IDs", h.config.AllowChannelIDs), zap.String("channel ID", i.ChannelID))
	if !lo.Contains(h.config.AllowChannelIDs, i.ChannelID) {
		h.logger.Debug("message hook in not allowed channel", zap.String("channelID", i.ChannelID))
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "このチャンネルからの呼出は許可されていません",
			},
		})
	}

	var totpClient TotpClient
	svcName := i.ApplicationCommandData().Options[0].Name
	tok, ok := h.config.Tokens.m[service(svcName)]
	if !ok {
		h.logger.Warn("svc not found", zap.String("service name", i.ApplicationCommandData().Name))
		return fmt.Errorf("not exist sub command")
	}
	totpClient = &TotpGen{secret: string(tok)}
	code, err := totpClient.GenerateCode(time.Now())
	if err != nil {
		h.logger.Error("generate code error", zap.Error(err), zap.String("serviceName", svcName))
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

type TotpGen struct {
	secret string
}

func (g *TotpGen) GenerateCode(t time.Time) (string, error) {

	return totp.GenerateCode(trimInnerWhite(g.secret), t)
}
func trimInnerWhite(secret string) string {
	return strings.Join(strings.Fields(secret), ``)
}

func IntractionCreateHandlerRouter(logger *zap.Logger, config Config, cmdName string) (func(s *discordgo.Session, i *discordgo.InteractionCreate) error, error) {
	switch cmdName {
	case "2fa":
		handler := &TotpHandler{logger: logger, config: config}
		return handler.HandleIntractionCreate, nil
	default:
		logger.Error("unknown interaction create", zap.String("cmdName", cmdName))
		return nil, fmt.Errorf("no route to handler")
	}
}
