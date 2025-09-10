package user

import (
	"context"

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
	userRepository domain.UserRepositoryInterface
}

func NewRegisterUseCase(userRepository domain.UserRepositoryInterface) *RegisterUseCase {
	return &RegisterUseCase{
		userRepository: userRepository,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, r RegisterRequest) (RegisterResponse, error) {
	user, err := domain.NewUser(r.Name, r.Email, r.Password)
	if err != nil {
		return RegisterResponse{}, err
	}

	if err = uc.userRepository.Create(ctx, user); err != nil {
		return RegisterResponse{}, err
	}

	return RegisterResponse{
		Id: user.Id,
	}, nil
}
