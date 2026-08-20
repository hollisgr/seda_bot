package handler

import (
	"context"
	"fmt"
	"log"
	"sedabot/internal/model"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) SaveEventType(ctx context.Context, b *bot.Bot, update *models.Update, user model.User) {
	chatId := update.Message.Chat.ID
	updType := update.Message.Text

	h.mu.Lock()
	draft := h.eventDrafts[user.TgId]
	draft.Type = updType
	h.eventDrafts[user.TgId] = draft
	h.mu.Unlock()

	err := h.userUC.SetState(ctx, user.TgId, model.EventNameAwaiting)
	if err != nil {
		log.Println("handler save event type err: ", err)
		h.sendMessage(ctx, b, chatId, "internal error, try later")
		return
	}

	h.sendMessage(ctx, b, chatId, "enter event name")
}

func (h *Handler) SaveEventName(ctx context.Context, b *bot.Bot, update *models.Update, user model.User) {
	chatId := update.Message.Chat.ID
	updName := update.Message.Text

	h.mu.Lock()
	draft := h.eventDrafts[user.TgId]
	draft.Name = updName
	h.eventDrafts[user.TgId] = draft
	h.mu.Unlock()

	err := h.userUC.SetState(ctx, user.TgId, model.EventDescriptionAwaiting)
	if err != nil {
		log.Println("handler save event name err: ", err)
		h.sendMessage(ctx, b, chatId, "internal error, try later")
		return
	}

	h.sendMessage(ctx, b, chatId, "enter event description")
}

func (h *Handler) SaveEventDescription(ctx context.Context, b *bot.Bot, update *models.Update, user model.User) {
	chatId := update.Message.Chat.ID
	updDesc := update.Message.Text

	h.mu.Lock()
	draft := h.eventDrafts[user.TgId]
	draft.Description = updDesc
	h.eventDrafts[user.TgId] = draft
	h.mu.Unlock()

	err := h.userUC.SetState(ctx, user.TgId, model.EventDateAwaiting)
	if err != nil {
		log.Println("handler save event description err: ", err)
		h.sendMessage(ctx, b, chatId, "internal error, try later")
		return
	}

	h.sendMessage(ctx, b, chatId, "enter event date (layout: dd-mm-yyyy)")
}

func (h *Handler) SaveEventDate(ctx context.Context, b *bot.Bot, update *models.Update, user model.User) {
	chatId := update.Message.Chat.ID
	updDate := update.Message.Text
	layout := "01-02-2006"
	date, err := time.Parse(layout, updDate)
	if err != nil {
		log.Println("handler save event date err: ", err)
		h.sendMessage(ctx, b, chatId, "internal error, try later")
		return
	}

	h.mu.Lock()
	draft := h.eventDrafts[user.TgId]
	draft.Date = date
	h.eventDrafts[user.TgId] = draft
	h.mu.Unlock()

	err = h.userUC.SetState(ctx, user.TgId, model.EventTimeAwaiting)
	if err != nil {
		log.Println("handler save event date err: ", err)
		h.sendMessage(ctx, b, chatId, "internal error, try later")
		return
	}

	h.sendMessage(ctx, b, chatId, "enter event time (layout: hh:mm)")
}

func (h *Handler) SaveEventTime(ctx context.Context, b *bot.Bot, update *models.Update, user model.User) {
	chatId := update.Message.Chat.ID
	updTime := update.Message.Text
	layout := "15:04"
	eventTime, err := time.Parse(layout, updTime)
	if err != nil {
		log.Println("handler save event time err: ", err)
		h.sendMessage(ctx, b, chatId, "internal error, try later")
		return
	}

	h.mu.Lock()
	draft, exists := h.eventDrafts[user.TgId]
	if !exists {
		h.mu.Unlock()
		h.sendMessage(ctx, b, chatId, "draft not found")
		return
	}

	h.mu.Unlock()

	finalDate := time.Date(
		draft.Date.Year(),
		draft.Date.Month(),
		draft.Date.Day(),
		eventTime.Hour(),
		eventTime.Minute(),
		0,
		0,
		draft.Date.Location(),
	)

	draft.Date = finalDate

	id, err := h.eventUC.Save(ctx, draft)
	if err != nil {
		log.Println("handler save event time err: ", err)
		h.sendMessage(ctx, b, chatId, "internal error, try later")
		return
	}

	err = h.userUC.SetState(ctx, user.TgId, model.MainMenu)
	if err != nil {
		log.Println("handler save event time err: ", err)
		h.sendMessage(ctx, b, chatId, "internal error, try later")
		return
	}

	h.mu.Lock()
	delete(h.eventDrafts, user.TgId)
	h.mu.Unlock()

	text := fmt.Sprintf("event created! id: %d", id)

	h.sendMessageWithKeyboard(ctx, b, chatId, text, h.getMarkup(user.Role))
}

func (h *Handler) SaveEvent(ctx context.Context, b *bot.Bot, update *models.Update) {}

func (h *Handler) LoadEvent(ctx context.Context, b *bot.Bot, update *models.Update) {}

func (h *Handler) LoadActualEvents(ctx context.Context, b *bot.Bot, update *models.Update) {}
