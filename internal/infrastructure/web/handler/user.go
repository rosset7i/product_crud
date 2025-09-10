package handler

import (
	"net/http"

	"github.com/rosset7i/product_crud/internal/infrastructure/web"
	"github.com/rosset7i/product_crud/internal/usecase/user"
)

type UserHandler struct {
	registerUseCase *user.RegisterUseCase
	loginUseCase    *user.LoginUseCase
}

func NewUserHandler(registerUseCase *user.RegisterUseCase, loginUseCase *user.LoginUseCase) *UserHandler {
	return &UserHandler{
		registerUseCase: registerUseCase,
		loginUseCase:    loginUseCase,
	}
}

// Register godoc
// @Tags         Users
// @Param        request  body      user.RegisterRequest  true "payload"
// @Success      201      {object}  user.RegisterResponse
// @Failure      400      {object}  web.errorResponse
// @Router       /v1/users/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeJSONBody[user.RegisterRequest](r)
	if err != nil {
		web.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.registerUseCase.Execute(r.Context(), req)
	if err != nil {
		web.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	web.WriteJSON(w, http.StatusCreated, response)
}

// Login godoc
// @Tags         Users
// @Param        request  body      user.LoginRequest   true "payload"
// @Success      200      {object}  user.LoginResponse
// @Failure      400      {object}  web.errorResponse
// @Failure      401      {object}  web.errorResponse
// @Router       /v1/users/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeJSONBody[user.LoginRequest](r)
	if err != nil {
		web.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.loginUseCase.Execute(r.Context(), req)
	if err != nil {
		web.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	web.WriteJSON(w, http.StatusOK, response)
}
