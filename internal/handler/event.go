package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) SaveEvent(ctx context.Context, b *bot.Bot, update *models.Update)

func (h *Handler) LoadEvent(ctx context.Context, b *bot.Bot, update *models.Update)

func (h *Handler) LoadActualEvents(ctx context.Context, b *bot.Bot, update *models.Update)
