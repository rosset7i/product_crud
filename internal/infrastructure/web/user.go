package web

import (
	"encoding/json"
	"errors"
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

var (
	errInvalidUser = errors.New("invalid email or password")
)

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user with name, email, and password
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      user.RegisterRequest  true "RegisterRequest"
// @Success      201      {object}  user.RegisterResponse
// @Failure      400      {object}  errorResponse
// @Router       /users/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req user.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := user.NewRegisterUseCase(h.userRepository).Execute(req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, errInvalidUser.Error())
		return
	}

	writeJSON(w, http.StatusCreated, response)
}

// Login godoc
// @Summary      Authenticate a user
// @Description  Validates user credentials and returns a JWT token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      user.LoginRequest   true "LoginRequest"
// @Success      200      {object}  user.LoginResponse
// @Failure      400      {object}  errorResponse
// @Failure      401      {object}  errorResponse
// @Router       /users/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req user.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := user.NewLoginUseCase(h.userRepository, h.jwtAuth, h.jwtExpiresIn).Execute(req)
	if err != nil {
		writeError(w, http.StatusUnauthorized, errInvalidUser.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)
}
