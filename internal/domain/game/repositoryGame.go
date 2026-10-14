package game

import "context"

type GameAttemptRepository interface {
	// Submit guarda el intento junto con todas sus notas de forma atómica.
	Submit(ctx context.Context, attempt *GameAttempt) error
}

type LeaderboardEntry struct {
	UserID       string
	Username     string
	BestScore    int
	BestAccuracy float64
}

type LeaderboardRepository interface {
	GetByLevel(ctx context.Context, levelID string) ([]LeaderboardEntry, error)
}
