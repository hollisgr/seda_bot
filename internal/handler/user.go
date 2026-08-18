package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sedabot/internal/model"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UserUseCase interface {
	SaveUser(ctx context.Context, user model.User) (int, error)
	LoadUser(ctx context.Context, tgId int64) (model.User, error)
	LoadUserList(ctx context.Context, offset int, limit int) ([]model.User, error)
	SetRole(ctx context.Context, tgId int64, role model.Role) error
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

	err = h.userUC.SetRole(ctx, int64(tgId), roleInput)
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
