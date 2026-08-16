package usecase

type EventRepository interface {
}

type EventUseCase struct {
	eventRepo EventRepository
}

func NewEventUseCase(eventRepo EventRepository) *EventUseCase {
	return &EventUseCase{
		eventRepo: eventRepo,
	}
}
