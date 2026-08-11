package model

type User struct {
	Id        int    `db:"id"`
	TgId      int    `db:"tg_id"`
	ChatId    int    `db:"chat_id"`
	Name      string `db:"name"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}
