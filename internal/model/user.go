package model

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type State string

const (
	MainMenu                 State = "main_menu"
	EventTypeAwaiting        State = "event:type_awaiting"
	EventNameAwaiting        State = "event:name_awaiting"
	EventDescriptionAwaiting State = "event:description_awaiting"
	EventDateAwaiting        State = "event:date_awaiting"
	EventTimeAwaiting        State = "event:time_awaiting"
)

func (s State) IsValid() bool {
	switch s {
	case MainMenu,
		EventTypeAwaiting,
		EventNameAwaiting,
		EventDescriptionAwaiting,
		EventDateAwaiting,
		EventTimeAwaiting:
		return true
	default:
		return false
	}
}

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
