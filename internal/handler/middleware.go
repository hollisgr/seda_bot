package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) AdminOnly(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		if update.Message == nil || update.Message.From == nil {
			return
		}

		if update.Message.From.ID != h.cfg.OwnerId {
			h.sendMessage(ctx, bot, update.Message.Chat.ID, "access denied")
			return
		}

		next(ctx, bot, update)
	}
}
