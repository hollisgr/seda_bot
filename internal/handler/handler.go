package handler

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UserUseCase interface{}

type Handler struct {
	userUC UserUseCase
}

func New(user UserUseCase) *Handler {
	return &Handler{
		userUC: user,
	}
}

func (h *Handler) Register(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.StartHandler)
}

func (h *Handler) Default(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}

func (h *Handler) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	firstName := update.Message.From.FirstName
	text := fmt.Sprintf("Hello, %s!\nWelcome to SedaBot!", firstName)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}
