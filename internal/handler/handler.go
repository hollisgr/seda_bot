package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sedabot/internal/config"
	"sedabot/internal/model"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type EventUseCase interface {
	Save(ctx context.Context, event model.Event) (int, error)
	Load(ctx context.Context, id int) (model.Event, error)
	LoadActive(ctx context.Context) ([]model.Event, error)
}

type Handler struct {
	cfg           *config.Config
	userUC        UserUseCase
	eventUC       EventUseCase
	adminKeyBoard *models.ReplyKeyboardMarkup
	userKeyBoard  *models.ReplyKeyboardMarkup
	eventDrafts   map[int64]model.Event
	mu            sync.Mutex
}

func New(user UserUseCase, event EventUseCase, cfg *config.Config) *Handler {
	return &Handler{
		userUC:        user,
		eventUC:       event,
		cfg:           cfg,
		adminKeyBoard: newAdminKeyboard(),
		userKeyBoard:  newUserKeyboard(),
		eventDrafts:   make(map[int64]model.Event),
		mu:            sync.Mutex{},
	}
}

func (h *Handler) Register(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.StartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/user_list", bot.MatchTypeExact, h.OwnerOnly(h.UserList))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/set_role", bot.MatchTypePrefix, h.OwnerOnly(h.SetRole))
}

func (h *Handler) Default(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	fromUser := update.Message.From
	chatId := update.Message.Chat.ID

	user, err := h.userUC.LoadUser(ctx, fromUser.ID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			h.sendMessage(ctx, b, chatId, "user not found, try /start")
			return
		}
		h.sendMessage(ctx, b, chatId, "internal error, try later")
		return
	}

	if user.Role != model.RoleAdmin {
		h.sendMessageWithKeyboard(ctx, b, chatId, "main menu", h.getMarkup(user.Role))
	}

	switch user.State {
	case model.EventTypeAwaiting:
		h.SaveEventType(ctx, b, update, user)
	case model.EventNameAwaiting:
		h.SaveEventName(ctx, b, update, user)
	case model.EventDescriptionAwaiting:
		h.SaveEventDescription(ctx, b, update, user)
	case model.EventDateAwaiting:
		h.SaveEventDate(ctx, b, update, user)
	case model.EventTimeAwaiting:
		h.SaveEventTime(ctx, b, update, user)
	default:
		h.sendMessageWithKeyboard(ctx, b, chatId, "main menu", h.getMarkup(user.Role))
	}
}

func (h *Handler) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	chatId := update.Message.Chat.ID
	fromUser := update.Message.From

	var text string
	var user model.User

	user, err := h.userUC.LoadUser(ctx, fromUser.ID)

	switch {
	case err == nil:
		text = fmt.Sprintf("Hello, %s!\nWelcome back!", fromUser.FirstName)

	case errors.Is(err, model.ErrNotFound):
		newUser := model.User{
			TgId:      fromUser.ID,
			ChatId:    chatId,
			Name:      fromUser.Username,
			FirstName: fromUser.FirstName,
			LastName:  fromUser.LastName,
			Role:      model.RoleUser,
		}
		id, err := h.userUC.SaveUser(ctx, newUser)

		if err != nil {
			log.Println("start handler save user err: ", err)
			text = "internal error, try later"
		} else {
			text = fmt.Sprintf("Hello, %s!\nWelcome to SedaBot!", fromUser.FirstName)
			log.Println("start handler save user success, id: ", id)
		}

		user = newUser

	default:
		log.Println("start handler load user err: ", err)
		text = "internal error, try later"
	}

	h.sendMessageWithKeyboard(ctx, b, chatId, text, h.getMarkup(user.Role))
}

func (h *Handler) sendMessage(ctx context.Context, b *bot.Bot, chatId int64, text string) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   text,
	})
}

func (h *Handler) sendMessageWithKeyboard(ctx context.Context, b *bot.Bot, chatId int64, text string, markup *models.ReplyKeyboardMarkup) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        text,
		ReplyMarkup: markup,
	})
}

func newUserKeyboard() *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "USER BTN 1"},
				{Text: "USER BTN 2"},
			},
			{
				{Text: "USER BTN 3"},
				{Text: "USER BTN 4"},
			},
		},
		ResizeKeyboard: true,
	}
}

func newAdminKeyboard() *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "ADMIN BTN 1"},
				{Text: "ADMIN BTN 2"},
			},
			{
				{Text: "ADMIN BTN 3"},
				{Text: "ADMIN BTN 4"},
			},
		},
		ResizeKeyboard: true,
	}
}

func (h *Handler) getMarkup(role model.Role) *models.ReplyKeyboardMarkup {
	if role == model.RoleAdmin {
		return h.adminKeyBoard
	}
	return h.userKeyBoard
}
