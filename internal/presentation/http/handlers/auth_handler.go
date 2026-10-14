package handlers

import (
	"encoding/json"
	"net/http"

	"mg-back/internal/application/auth"
	"mg-back/internal/presentation/http/response"
	"mg-back/internal/presentation/http/validation"
)

type AuthHandler struct {
	registerUser *auth.RegisterUserUseCase
	login        *auth.LoginUseCase
	refreshToken *auth.RefreshTokenUseCase
}

func NewAuthHandler(
	registerUser *auth.RegisterUserUseCase,
	login *auth.LoginUseCase,
	refreshToken *auth.RefreshTokenUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUser: registerUser,
		login:        login,
		refreshToken: refreshToken,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input auth.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "JSON malformado"})
		return
	}

	if errs := validation.Validate(input); errs != nil {
		response.JSON(w, http.StatusUnprocessableEntity, map[string]interface{}{"errors": errs})
		return
	}

	result, err := h.registerUser.Execute(r.Context(), input)
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, result)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input auth.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "JSON malformado"})
		return
	}

	if errs := validation.Validate(input); errs != nil {
		response.JSON(w, http.StatusUnprocessableEntity, map[string]interface{}{"errors": errs})
		return
	}

	result, err := h.login.Execute(r.Context(), input)
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input auth.RefreshTokenInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "JSON malformado"})
		return
	}

	if errs := validation.Validate(input); errs != nil {
		response.JSON(w, http.StatusUnprocessableEntity, map[string]interface{}{"errors": errs})
		return
	}

	result, err := h.refreshToken.Execute(r.Context(), input.RefreshToken)
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
