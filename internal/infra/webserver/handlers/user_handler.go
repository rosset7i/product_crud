package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/rosset7i/product_crud/config"
	"github.com/rosset7i/product_crud/internal/dto"
	"github.com/rosset7i/product_crud/internal/entity"
	"github.com/rosset7i/product_crud/internal/infra/database"
	"github.com/rosset7i/product_crud/internal/infra/webserver"
)

type UserHandler struct {
	UserDB       *database.UserRepository
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn time.Duration
}

func NewUserHandler(userDb *database.UserRepository, c *config.Conf) *UserHandler {
	return &UserHandler{
		UserDB:       userDb,
		Jwt:          c.Auth.JwtAuth,
		JwtExpiresIn: c.Auth.JwtExpiresIn,
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
// @Param        request  body      dto.CreateUserRequest  true  "User registration request"
// @Success      201      {object}  entity.User
// @Failure      400      {object}  webserver.ErrorResponse
// @Failure      500      {object}  webserver.ErrorResponse
// @Router       /users/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		webserver.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := entity.NewUser(req.Name, req.Email, req.Password)
	if err != nil {
		webserver.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.UserDB.Create(user); err != nil {
		webserver.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	webserver.WriteJSON(w, http.StatusCreated, user)
}

// Login godoc
// @Summary      Authenticate a user
// @Description  Validates user credentials and returns a JWT token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest   true  "User login request"
// @Success      200      {object}  dto.LoginResponse
// @Failure      400      {object}  webserver.ErrorResponse
// @Failure      401      {object}  webserver.ErrorResponse
// @Failure      404      {object}  webserver.ErrorResponse
// @Failure      500      {object}  webserver.ErrorResponse
// @Router       /users/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		webserver.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.UserDB.FetchByEmail(req.Email)
	if err != nil {
		webserver.WriteError(w, http.StatusNotFound, errInvalidUser.Error())
		return
	}

	if !user.ValidatePassword(req.Password) {
		webserver.WriteError(w, http.StatusUnauthorized, errInvalidUser.Error())
		return
	}

	_, tokenString, err := h.Jwt.Encode(map[string]any{
		"sub": user.Id.String(),
		"exp": time.Now().Add(h.JwtExpiresIn).Unix(),
	})
	if err != nil {
		webserver.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	webserver.WriteJSON(w, http.StatusOK, dto.LoginResponse{AccessToken: tokenString})
}
