package user

import (
	"context"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/rosset7i/product_crud/internal/domain"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type LoginUseCase struct {
	userRepository domain.UserRepositoryInterface
	jwtAuth        *jwtauth.JWTAuth
	jtwExpiresIn   time.Duration
}

func NewLoginUseCase(userRepository domain.UserRepositoryInterface, jwtAuth *jwtauth.JWTAuth, jtwExpiresIn time.Duration) *LoginUseCase {
	return &LoginUseCase{
		userRepository: userRepository,
		jwtAuth:        jwtAuth,
		jtwExpiresIn:   jtwExpiresIn,
	}
}

func (uc *LoginUseCase) Execute(r LoginRequest) (LoginResponse, error) {
	user, err := uc.userRepository.FetchByEmail(context.TODO(), r.Email)
	if err != nil {
		return LoginResponse{}, err
	}

	if !user.ValidatePassword(r.Password) {
		return LoginResponse{}, err
	}

	_, tokenString, err := uc.jwtAuth.Encode(map[string]any{
		"sub": user.Id.String(),
		"exp": time.Now().Add(uc.jtwExpiresIn).Unix(),
	})
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken: tokenString,
	}, nil
}
