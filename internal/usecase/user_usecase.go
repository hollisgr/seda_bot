package usecase

import (
	"context"
	"fmt"
	"sedabot/internal/model"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user model.User) (int, error)
	LoadUserByTgId(ctx context.Context, tgId int64) (model.User, error)
	LoadUserList(ctx context.Context, offset int, limit int) ([]model.User, error)
	SetRole(ctx context.Context, tgId int64, role model.Role) error
}

type UserUseCase struct {
	userRepo UserRepository
	cache    *expirable.LRU[int64, model.User]
}

func NewUserUseCase(userRepo UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		cache:    expirable.NewLRU[int64, model.User](300, nil, time.Minute*30),
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

	u.cache.Remove(tgId)

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
	cashedUser, ok := u.getUserFromCache(tgId)
	if ok {
		return cashedUser, nil
	}

	user, err := u.userRepo.LoadUserByTgId(ctx, tgId)
	if err != nil {
		return model.User{}, err
	}

	u.cache.Add(user.TgId, user)
	return user, nil
}

func (u *UserUseCase) getUserFromCache(tgId int64) (model.User, bool) {
	user, ok := u.cache.Get(tgId)
	if !ok {
		return model.User{}, false
	}

	return user, true
}
