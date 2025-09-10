package web

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/rosset7i/product_crud/config"
	"github.com/rosset7i/product_crud/internal/infrastructure/database"
	"github.com/rosset7i/product_crud/internal/usecase/user"
)

type UserHandler struct {
	userRepository *database.UserRepository
	jwtAuth        *jwtauth.JWTAuth
	jwtExpiresIn   time.Duration
}

func NewUserHandler(userRepository *database.UserRepository, c *config.Conf) *UserHandler {
	return &UserHandler{
		userRepository: userRepository,
		jwtAuth:        c.Auth.JwtAuth,
		jwtExpiresIn:   c.Auth.JwtExpiresIn,
	}
}

// Register godoc
// @Tags         Users
// @Param        request  body      user.RegisterRequest  true "payload"
// @Success      201      {object}  user.RegisterResponse
// @Failure      400      {object}  errorResponse
// @Router       /v1/users/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	req, err := decodeJSONBody[user.RegisterRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := user.NewRegisterUseCase(h.userRepository).Execute(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, response)
}

// Login godoc
// @Tags         Users
// @Param        request  body      user.LoginRequest   true "payload"
// @Success      200      {object}  user.LoginResponse
// @Failure      400      {object}  errorResponse
// @Failure      401      {object}  errorResponse
// @Router       /v1/users/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	req, err := decodeJSONBody[user.LoginRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := user.NewLoginUseCase(h.userRepository, h.jwtAuth, h.jwtExpiresIn).Execute(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)
}
