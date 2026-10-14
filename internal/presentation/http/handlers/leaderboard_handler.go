package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"mg-back/internal/application/leaderboard"
	"mg-back/internal/presentation/http/response"
)

type LeaderboardHandler struct {
	getLeaderboard *leaderboard.GetLeaderboardUseCase
}

func NewLeaderboardHandler(getLeaderboard *leaderboard.GetLeaderboardUseCase) *LeaderboardHandler {
	return &LeaderboardHandler{getLeaderboard: getLeaderboard}
}

func (h *LeaderboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	result, err := h.getLeaderboard.Execute(r.Context(), chi.URLParam(r, "levelID"))
	if err != nil {
		response.FromDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}
