package postgres

import (
	"context"
	"errors"
	"fmt"
	"sedabot/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) SaveUser(ctx context.Context, user model.User) (int, error) {
	var id int
	query := `
		INSERT INTO users(tg_id, chat_id, first_name, last_name, name)
		VALUES (@tg_id, @chat_id, @first_name, @last_name, @name)
		RETURNING id
	`
	args := pgx.NamedArgs{
		"tg_id":      user.TgId,
		"chat_id":    user.ChatId,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"name":       user.Name,
	}

	row := s.db.QueryRow(ctx, query, args)
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("db save user error: %w", err)
	}

	return id, nil
}

func (s *Storage) LoadUserByTgId(ctx context.Context, tgId int) (model.User, error) {
	var res model.User
	query := `
		SELECT
			id,
			tg_id, 
			chat_id, 
			first_name, 
			last_name, 
			name,
			role
		FROM 
			users
		WHERE
			tg_id = @tg_id
	`
	args := pgx.NamedArgs{
		"tg_id": tgId,
	}

	rows, err := s.db.Query(ctx, query, args)

	if err != nil {
		return res, fmt.Errorf("db load user by tg_id query error: %w", err)
	}
	defer rows.Close()

	res, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return res, model.ErrNotFound
		}
		return res, fmt.Errorf("db load user by tg_id collect error: %w", err)
	}

	return res, nil
}

func (s *Storage) LoadUserList(ctx context.Context, offset int, limit int) ([]model.User, error) {
	var res []model.User

	query := `
		SELECT
			id,
			tg_id,
			chat_id,
			name,
			first_name,
			last_name,
			role
		FROM users
		LIMIT @limit
		OFFSET @offset
	`

	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := s.db.Query(ctx, query, args)
	if err != nil {
		return res, fmt.Errorf("db load user list query err: %w", err)
	}

	defer rows.Close()

	res, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return res, model.ErrNotFound
		}
		return res, fmt.Errorf("db load user list collect err: %w", err)
	}
	return res, nil
}

func (s *Storage) SetRole(ctx context.Context, tgId int, role model.Role) error {
	var id int
	query := `
		UPDATE users
		SET
			role = @role
		WHERE
			tg_id = @tg_id
		RETURNING id
	`

	args := pgx.NamedArgs{
		"role":  role,
		"tg_id": tgId,
	}

	row := s.db.QueryRow(ctx, query, args)
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrNotFound
		}
		return fmt.Errorf("db set role row scan error: %w", err)
	}
	return nil
}
