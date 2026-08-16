package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type EventStorage struct {
	db *pgxpool.Pool
}

func NewEventStorage(p *pgxpool.Pool) *EventStorage {
	return &EventStorage{
		db: p,
	}
}
