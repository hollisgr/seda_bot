package usecase

import (
	"context"
	"sedabot/internal/model"
)

type EventRepository interface {
	Save(ctx context.Context, event model.Event) (int, error)
	Load(ctx context.Context, id int) (model.Event, error)
	LoadActiveEvents(ctx context.Context) ([]model.Event, error)
}

type EventUseCase struct {
	eventRepo EventRepository
}

func NewEventUseCase(eventRepo EventRepository) *EventUseCase {
	return &EventUseCase{
		eventRepo: eventRepo,
	}
}

func (u *EventUseCase) Save(ctx context.Context, event model.Event) (int, error) {
	id, err := u.eventRepo.Save(ctx, event)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *EventUseCase) Load(ctx context.Context, id int) (model.Event, error) {
	event, err := u.eventRepo.Load(ctx, id)
	if err != nil {
		return model.Event{}, err
	}
	return event, nil
}

func (u *EventUseCase) LoadActive(ctx context.Context) ([]model.Event, error) {
	events, err := u.eventRepo.LoadActiveEvents(ctx)
	if err != nil {
		return nil, err
	}
	return events, nil
}
