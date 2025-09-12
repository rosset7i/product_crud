package user

import (
	"context"
	"errors"
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
	userRepository domain.UserRepository
	jwtAuth        *jwtauth.JWTAuth
	jtwExpiresIn   time.Duration
}

func NewLoginUseCase(userRepository domain.UserRepository, jwtAuth *jwtauth.JWTAuth, jtwExpiresIn time.Duration) *LoginUseCase {
	return &LoginUseCase{
		userRepository: userRepository,
		jwtAuth:        jwtAuth,
		jtwExpiresIn:   jtwExpiresIn,
	}
}

var (
	errWrongEmailOrPassword  = errors.New("wrong email or password")
	errCouldNotGenerateToken = errors.New("could not generate token, try again")
)

func (uc *LoginUseCase) Execute(ctx context.Context, r LoginRequest) (LoginResponse, error) {
	user, err := uc.userRepository.FetchByEmail(ctx, r.Email)
	if err != nil {
		return LoginResponse{}, errWrongEmailOrPassword
	}

	if !user.ValidatePassword(r.Password) {
		return LoginResponse{}, errWrongEmailOrPassword
	}

	_, tokenString, err := uc.jwtAuth.Encode(map[string]any{
		"sub": user.Id.String(),
		"exp": time.Now().Add(uc.jtwExpiresIn).Unix(),
	})
	if err != nil {
		return LoginResponse{}, errCouldNotGenerateToken
	}

	return LoginResponse{
		AccessToken: tokenString,
	}, nil
}
