package handler

import (
	"context"
	"sedabot/internal/model"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) OwnerOnly(next bot.HandlerFunc) bot.HandlerFunc {
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

func (h *Handler) AdminOnly(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		if update.Message == nil || update.Message.From == nil {
			return
		}

		from := update.Message.From

		user, err := h.userUC.LoadUser(ctx, from.ID)
		if err != nil {
			h.sendMessage(ctx, bot, update.Message.Chat.ID, "access denied")
			return
		}

		if user.Role != model.RoleAdmin {
			h.sendMessage(ctx, bot, update.Message.Chat.ID, "access denied")
			return
		}

		next(ctx, bot, update)
	}
}
