package model

import "time"

type Event struct {
	Id          int       `db:"id"`
	Type        string    `db:"type"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Date        time.Time `db:"date"`
}
