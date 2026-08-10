package usecase

type UserRepository interface{}

type UserUseCase struct {
	UserRepo UserRepository
}

func NewUserUseCase(userRepo UserRepository) *UserUseCase {
	return &UserUseCase{
		UserRepo: userRepo,
	}
}
