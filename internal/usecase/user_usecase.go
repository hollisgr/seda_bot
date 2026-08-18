package usecase

import (
	"context"
	"fmt"
	"sedabot/internal/model"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user model.User) (int, error)
	LoadUserByTgId(ctx context.Context, tgId int64) (model.User, error)
	LoadUserList(ctx context.Context, offset int, limit int) ([]model.User, error)
	SetRole(ctx context.Context, tgId int64, role model.Role) error
}

type UserUseCase struct {
	userRepo UserRepository
}

func NewUserUseCase(userRepo UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u *UserUseCase) LoadUserList(ctx context.Context, offset int, limit int) ([]model.User, error) {
	list, err := u.userRepo.LoadUserList(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (u *UserUseCase) SetRole(ctx context.Context, tgId int64, role model.Role) error {
	if role != model.RoleAdmin && role != model.RoleUser {
		return fmt.Errorf("user usecase set role err: invalid role %s", role)
	}

	err := u.userRepo.SetRole(ctx, tgId, role)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) SaveUser(ctx context.Context, user model.User) (int, error) {
	res, err := u.userRepo.SaveUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (u *UserUseCase) LoadUser(ctx context.Context, tgId int64) (model.User, error) {
	user, err := u.userRepo.LoadUserByTgId(ctx, tgId)
	if err != nil {
		return user, err
	}
	return user, nil
}
