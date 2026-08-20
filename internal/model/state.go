package model

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
