package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"mg-back/internal/domain/game"
)

type LeaderboardRepository struct {
	pool *pgxpool.Pool
}

func NewLeaderboardRepository(pool *pgxpool.Pool) *LeaderboardRepository {
	return &LeaderboardRepository{pool: pool}
}

func (r *LeaderboardRepository) GetByLevel(ctx context.Context, levelID string) ([]game.LeaderboardEntry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT l.user_id, u.username, l.best_score, l.best_accuracy
		FROM leaderboard l
		JOIN users u ON u.id = l.user_id
		WHERE l.level_id = $1
		ORDER BY l.best_score DESC, l.best_accuracy DESC
		LIMIT 100
	`, levelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []game.LeaderboardEntry
	for rows.Next() {
		var e game.LeaderboardEntry
		if err := rows.Scan(&e.UserID, &e.Username, &e.BestScore, &e.BestAccuracy); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
