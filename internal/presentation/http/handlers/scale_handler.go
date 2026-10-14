package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"mg-back/internal/application/scale"
	"mg-back/internal/presentation/http/response"
)

type ScaleHandler struct {
	getScales      *scale.GetScalesUseCase
	getLevel       *scale.GetLevelUseCase
	getScaleLevels *scale.GetScaleLevelsUseCase
}

func NewScaleHandler(
	getScales *scale.GetScalesUseCase,
	getLevel *scale.GetLevelUseCase,
	getScaleLevels *scale.GetScaleLevelsUseCase,
) *ScaleHandler {
	return &ScaleHandler{getScales: getScales, getLevel: getLevel, getScaleLevels: getScaleLevels}
}

func (h *ScaleHandler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.getScales.Execute(r.Context())
	if err != nil {
		response.FromDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *ScaleHandler) ListLevels(w http.ResponseWriter, r *http.Request) {
	result, err := h.getScaleLevels.Execute(r.Context(), chi.URLParam(r, "scaleID"))
	if err != nil {
		response.FromDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *ScaleHandler) GetLevel(w http.ResponseWriter, r *http.Request) {
	result, err := h.getLevel.Execute(r.Context(), chi.URLParam(r, "levelID"))
	if err != nil {
		response.FromDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}
