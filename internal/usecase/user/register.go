package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rosset7i/product_crud/internal/domain"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Id uuid.UUID `json:"id"`
}

type RegisterUseCase struct {
	userRepository domain.UserRepository
}

func NewRegisterUseCase(userRepository domain.UserRepository) *RegisterUseCase {
	return &RegisterUseCase{
		userRepository: userRepository,
	}
}

var (
	errUserAlreadyExists   = errors.New("user with that email already exists")
	errCouldNotPersistUser = errors.New("failed to persist user")
)

func (uc *RegisterUseCase) Execute(ctx context.Context, r RegisterRequest) (RegisterResponse, error) {
	user, _ := uc.userRepository.FetchByEmail(ctx, r.Email)
	if user != nil {
		return RegisterResponse{}, errUserAlreadyExists
	}

	newUser, err := domain.NewUser(r.Name, r.Email, r.Password)
	if err != nil {
		return RegisterResponse{}, err
	}

	if err = uc.userRepository.Create(ctx, newUser); err != nil {
		return RegisterResponse{}, errCouldNotPersistUser
	}

	return RegisterResponse{
		Id: newUser.Id,
	}, nil
}
