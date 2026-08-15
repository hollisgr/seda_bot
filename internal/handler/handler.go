package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sedabot/internal/config"
	"sedabot/internal/model"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UserUseCase interface {
	SaveUser(ctx context.Context, user model.User) (int, error)
	LoadUser(ctx context.Context, tgId int) (model.User, error)
	LoadUserList(ctx context.Context, offset int, limit int) ([]model.User, error)
	SetRole(ctx context.Context, tgId int, role model.Role) error
}

type Handler struct {
	cfg    *config.Config
	userUC UserUseCase
}

func New(user UserUseCase, cfg *config.Config) *Handler {
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

func (h *Handler) UserList(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}
	chatId := update.Message.Chat.ID
	list, err := h.userUC.LoadUserList(ctx, 0, 20)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			log.Println("userlist handler err: users not found")
			h.sendMessage(ctx, b, chatId, "users not found")
			return
		}
		log.Println("userlist handler err: ", err)
		h.sendMessage(ctx, b, chatId, "internal error")
		return
	}

	var builder strings.Builder
	builder.WriteString("<b>User List</b>\n\n")
	builder.WriteString("id, tg_id, username, first name, last name, role\n")

	for _, user := range list {
		line := fmt.Sprintf("%d) %d %s %s %s %s\n", user.Id, user.TgId, user.Name, user.FirstName, user.LastName, user.Role)
		builder.WriteString(line)
	}
	h.sendMessage(ctx, b, chatId, builder.String())
}

func (h *Handler) SetRole(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	chatId := update.Message.Chat.ID
	command := update.Message.Text
	parts := strings.Split(command, ":")

	if len(parts) != 3 {
		h.sendMessage(ctx, b, chatId, "incorrect command: use /set_role:role:tg_id")
		return
	}

	roleInput := model.Role(strings.TrimSpace(parts[1]))
	tgIdInput := strings.TrimSpace(parts[2])

	if roleInput != model.RoleAdmin && roleInput != model.RoleUser {
		h.sendMessage(ctx, b, chatId, "incorrect command: use /set_role:role:tg_id")
		return
	}

	tgId, err := strconv.Atoi(tgIdInput)
	if err != nil {
		log.Println("set role handler conv err: ", err)
		h.sendMessage(ctx, b, chatId, "internal error")
		return
	}

	err = h.userUC.SetRole(ctx, tgId, roleInput)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			h.sendMessage(ctx, b, chatId, "user not found")
			return
		}

		log.Println("set role handler err: ", err)
		h.sendMessage(ctx, b, chatId, "internal error")
		return
	}

	h.sendMessage(ctx, b, chatId, fmt.Sprintf("The user's role with tg_id: %d has been successfully changed to %s", tgId, roleInput))
}

func (h *Handler) sendMessage(ctx context.Context, b *bot.Bot, chatId int64, text string) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   text,
	})
}
