package leaderboard

import (
	"context"

	"mg-back/internal/domain/game"
)

type GetLeaderboardUseCase struct {
	repo game.LeaderboardRepository
}

func NewGetLeaderboardUseCase(repo game.LeaderboardRepository) *GetLeaderboardUseCase {
	return &GetLeaderboardUseCase{repo: repo}
}

func (uc *GetLeaderboardUseCase) Execute(ctx context.Context, levelID string) ([]Entry, error) {
	rows, err := uc.repo.GetByLevel(ctx, levelID)
	if err != nil {
		return nil, err
	}

	out := make([]Entry, 0, len(rows))
	for _, row := range rows {
		out = append(out, Entry{
			UserID:       row.UserID,
			Username:     row.Username,
			BestScore:    row.BestScore,
			BestAccuracy: row.BestAccuracy,
		})
	}
	return out, nil
}
