package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sedabot/internal/config"
	"sedabot/internal/model"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UserUseCase interface {
	SaveUser(ctx context.Context, user model.User) (int, error)
	LoadUser(ctx context.Context, tgId int) (model.User, error)
	LoadUserList(ctx context.Context, offset int, limit int) ([]model.User, error)
	SetRole(ctx context.Context, tgId int, role model.Role) error
}

type EventUseCase interface{}

type Handler struct {
	cfg    *config.Config
	userUC UserUseCase
}

func New(user UserUseCase, event EventUseCase, cfg *config.Config) *Handler {
	return &Handler{
		userUC: user,
		cfg:    cfg,
	}
}

func (h *Handler) Register(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.StartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/user_list", bot.MatchTypeExact, h.AdminOnly(h.UserList))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/set_role", bot.MatchTypePrefix, h.AdminOnly(h.SetRole))
}

func (h *Handler) Default(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}
	h.sendMessage(ctx, b, update.Message.Chat.ID, update.Message.Text)
}

func (h *Handler) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	chatId := update.Message.Chat.ID
	fromUser := update.Message.From

	var text string

	_, err := h.userUC.LoadUser(ctx, int(fromUser.ID))

	switch {
	case err == nil:
		text = fmt.Sprintf("Hello, %s!\nWelcome back!", fromUser.FirstName)

	case errors.Is(err, model.ErrNotFound):
		id, err := h.userUC.SaveUser(ctx, model.User{
			TgId:      fromUser.ID,
			ChatId:    chatId,
			Name:      fromUser.Username,
			FirstName: fromUser.FirstName,
			LastName:  fromUser.LastName,
		})

		if err != nil {
			log.Println("start handler save user err: ", err)
			text = "internal error, try later"
		} else {
			text = fmt.Sprintf("Hello, %s!\nWelcome to SedaBot!", fromUser.FirstName)
			log.Println("start handler save user success, id: ", id)
		}

	default:
		log.Println("start handler load user err: ", err)
		text = "internal error, try later"
	}

	h.sendMessage(ctx, b, chatId, text)
}

func (h *Handler) sendMessage(ctx context.Context, b *bot.Bot, chatId int64, text string) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   text,
	})
}
