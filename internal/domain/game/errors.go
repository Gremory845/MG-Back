package game

import "mg-back/internal/domain/shared"

var (
	ErrGameAttemptNotFound = shared.NewDomainError("GAME_ATTEMPT_NOT_FOUND", "intento de juego no encontrado")
	ErrGameAttemptInvalid  = shared.NewDomainError("GAME_ATTEMPT_INVALID", "el resultado del intento no es válido")
)
