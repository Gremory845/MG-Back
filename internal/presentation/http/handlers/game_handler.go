package handlers

import (
	"encoding/json"
	"net/http"

	"mg-back/internal/application/game"
	"mg-back/internal/presentation/http/middleware"
	"mg-back/internal/presentation/http/response"
	"mg-back/internal/presentation/http/validation"
)

type GameHandler struct {
	submitAttempt *game.SubmitAttemptUseCase
}

func NewGameHandler(submitAttempt *game.SubmitAttemptUseCase) *GameHandler {
	return &GameHandler{submitAttempt: submitAttempt}
}

func (h *GameHandler) SubmitAttempt(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.JSON(w, http.StatusUnauthorized, map[string]interface{}{
			"error": map[string]string{"code": "UNAUTHORIZED", "message": "no autenticado"},
		})
		return
	}

	var input game.SubmitAttemptInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "JSON malformado"})
		return
	}

	if errs := validation.Validate(input); errs != nil {
		response.JSON(w, http.StatusUnprocessableEntity, map[string]interface{}{"errors": errs})
		return
	}

	result, err := h.submitAttempt.Execute(r.Context(), userID, input)
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, result)
}
