package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(p *pgxpool.Pool) *Storage {
	return &Storage{
		db: p,
	}
}
