package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Register(b *bot.Bot) {}

func (h *Handler) Default(ctx context.Context, b *bot.Bot, update *models.Update) {}
