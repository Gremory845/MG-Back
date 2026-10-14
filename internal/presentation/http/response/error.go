package response

import (
	"errors"
	"net/http"

	"mg-back/internal/domain/shared"
)

var statusByCode = map[string]int{
	"USER_NOT_FOUND":             http.StatusNotFound,
	"EMAIL_ALREADY_EXISTS":       http.StatusConflict,
	"USERNAME_TAKEN":             http.StatusConflict,
	"INVALID_CREDENTIALS":        http.StatusUnauthorized,
	"REFRESH_TOKEN_NOT_FOUND":    http.StatusUnauthorized,
	"REFRESH_TOKEN_EXPIRED":      http.StatusUnauthorized,
	"REFRESH_TOKEN_REUSED":       http.StatusUnauthorized,
	"SCALE_NOT_FOUND":            http.StatusNotFound,
	"SCALE_DEFINITION_NOT_FOUND": http.StatusNotFound,
	"LEVEL_NOT_FOUND":            http.StatusNotFound,
	"GAME_ATTEMPT_NOT_FOUND":     http.StatusNotFound,
	"GAME_ATTEMPT_INVALID":       http.StatusUnprocessableEntity,
}

func FromDomainError(w http.ResponseWriter, err error) {
	var domainErr *shared.DomainError
	if errors.As(err, &domainErr) {
		status, ok := statusByCode[domainErr.Code]
		if !ok {
			status = http.StatusUnprocessableEntity
		}
		JSON(w, status, map[string]interface{}{
			"error": map[string]string{
				"code":    domainErr.Code,
				"message": domainErr.Message,
			},
		})
		return
	}

	JSON(w, http.StatusInternalServerError, map[string]interface{}{
		"error": map[string]string{
			"code":    "INTERNAL_ERROR",
			"message": "ocurrió un error inesperado",
		},
	})
}
