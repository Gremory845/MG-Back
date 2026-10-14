package music

import "mg-back/internal/domain/shared"

var (
	ErrScaleNotFound           = shared.NewDomainError("SCALE_NOT_FOUND", "escala no encontrada")
	ErrScaleDefinitionNotFound = shared.NewDomainError("SCALE_DEFINITION_NOT_FOUND", "definición de escala no encontrada")
	ErrLevelNotFound           = shared.NewDomainError("LEVEL_NOT_FOUND", "nivel no encontrado")
)
