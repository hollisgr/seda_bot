package postgres

import (
	"context"
	"errors"
	"fmt"
	"sedabot/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventStorage struct {
	db *pgxpool.Pool
}

func NewEventStorage(p *pgxpool.Pool) *EventStorage {
	return &EventStorage{
		db: p,
	}
}

func (s *EventStorage) Save(ctx context.Context, event model.Event) (int, error) {
	var id int
	query := `
		INSERT INTO events (
			type,
			name,
			description,
			date
		)
		VALUES (
			@type,
			@name,
			@description,
			@date
		)
		RETURNING id
	`
	args := pgx.NamedArgs{
		"type":        event.Type,
		"name":        event.Name,
		"description": event.Description,
		"date":        event.Date,
	}

	row := s.db.QueryRow(ctx, query, args)
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("db save event err: %w", err)
	}

	return id, nil
}

func (s *EventStorage) Load(ctx context.Context, id int) (model.Event, error) {
	var res model.Event
	query := `
		SELECT 
			id,
			type,
			name,
			description,
			date
		FROM events
		WHERE id = @id
	`
	args := pgx.NamedArgs{
		"id": id,
	}

	row := s.db.QueryRow(ctx, query, args)
	err := row.Scan(&res)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Event{}, model.ErrNotFound
		}
		return model.Event{}, fmt.Errorf("db load event err: %w", err)
	}

	return res, nil
}

func (s *EventStorage) LoadActiveEvents(ctx context.Context) ([]model.Event, error) {
	query := `
		SELECT 
			id,
			type,
			name,
			description,
			date
		FROM events
		WHERE 
			date > NOW()
		ORDER BY id
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("db load active events err: %w", err)
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Event])
	if err != nil {
		return nil, fmt.Errorf("db load active events err: %w", err)
	}

	if len(res) == 0 {
		return nil, model.ErrNotFound
	}
	return res, nil
}
