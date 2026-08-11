package usecase

import (
	"context"
	"sedabot/internal/model"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user model.User) (int, error)
	LoadUserByTgId(ctx context.Context, tgId int) (model.User, error)
}

type UserUseCase struct {
	userRepo UserRepository
}

func NewUserUseCase(userRepo UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u *UserUseCase) SaveUser(ctx context.Context, user model.User) (int, error) {
	res, err := u.userRepo.SaveUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (u *UserUseCase) LoadUser(ctx context.Context, tgId int) (model.User, error) {
	user, err := u.userRepo.LoadUserByTgId(ctx, tgId)
	if err != nil {
		return user, err
	}
	return user, nil
}
