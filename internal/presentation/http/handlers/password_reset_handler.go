package handlers

import (
	"encoding/json"
	"net/http"

	"mg-back/internal/application/auth"
	"mg-back/internal/presentation/http/response"
	"mg-back/internal/presentation/http/validation"
)

type PasswordResetHandler struct {
	forgot *auth.ForgotPasswordUseCase
	verify *auth.VerifyResetCodeUseCase
	reset  *auth.ResetPasswordUseCase
}

func NewPasswordResetHandler(
	forgot *auth.ForgotPasswordUseCase,
	verify *auth.VerifyResetCodeUseCase,
	reset *auth.ResetPasswordUseCase,
) *PasswordResetHandler {
	return &PasswordResetHandler{forgot: forgot, verify: verify, reset: reset}
}

// decode lee y valida el JSON; devuelve false si ya respondió con error.
func decode(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "JSON malformado"})
		return false
	}
	if errs := validation.Validate(dst); errs != nil {
		response.JSON(w, http.StatusUnprocessableEntity, map[string]interface{}{"errors": errs})
		return false
	}
	return true
}

func (h *PasswordResetHandler) Forgot(w http.ResponseWriter, r *http.Request) {
	var input auth.ForgotPasswordInput
	if !decode(w, r, &input) {
		return
	}
	if err := h.forgot.Execute(r.Context(), input); err != nil {
		response.FromDomainError(w, err)
		return
	}
	// Misma respuesta exista o no el correo
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Si el correo está registrado, te enviamos un código de 6 dígitos.",
	})
}

func (h *PasswordResetHandler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	var input auth.VerifyResetCodeInput
	if !decode(w, r, &input) {
		return
	}
	result, err := h.verify.Execute(r.Context(), input)
	if err != nil {
		response.FromDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *PasswordResetHandler) Reset(w http.ResponseWriter, r *http.Request) {
	var input auth.ResetPasswordInput
	if !decode(w, r, &input) {
		return
	}
	if err := h.reset.Execute(r.Context(), input); err != nil {
		response.FromDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Contraseña actualizada. Inicia sesión con la nueva."})
}
