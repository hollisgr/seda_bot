package model

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type State string

const (
	MainMenu State = "main_menu"
)

type User struct {
	Id        int    `db:"id"`
	TgId      int64  `db:"tg_id"`
	ChatId    int64  `db:"chat_id"`
	Name      string `db:"name"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Role      Role   `db:"role"`
	State     State  `db:"state"`
}
