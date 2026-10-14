package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"mg-back/internal/domain/game"
)

type GameAttemptRepository struct {
	pool *pgxpool.Pool
}

func NewGameAttemptRepository(pool *pgxpool.Pool) *GameAttemptRepository {
	return &GameAttemptRepository{pool: pool}
}

// Submit guarda el GameAttempt junto con todos sus NoteAttempt en una
// sola transacción: o se guardan los dos, o no se guarda ninguno.
func (r *GameAttemptRepository) Submit(ctx context.Context, attempt *game.GameAttempt) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // no-op si el commit ya ocurrió

	_, err = tx.Exec(ctx, `
		INSERT INTO game_attempts
			(id, user_id, level_id, score, accuracy, mistakes, completed, duration_ms, started_at, finished_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, attempt.ID, attempt.UserID, attempt.LevelID, attempt.Score, attempt.Accuracy,
		attempt.Mistakes, attempt.Completed, attempt.DurationMs, attempt.StartedAt, attempt.FinishedAt)
	if err != nil {
		return err
	}

	for _, n := range attempt.Notes {
		_, err = tx.Exec(ctx, `
			INSERT INTO note_attempts
				(game_attempt_id, position, target_pitch_class, target_octave,
				 detected_pitch_class, detected_octave, target_frequency, detected_frequency,
				 cents_difference, accuracy, result)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`, attempt.ID, n.Position, n.TargetPitchClass, n.TargetOctave,
			n.DetectedPitchClass, n.DetectedOctave, n.TargetFrequency, n.DetectedFrequency,
			n.CentsDifference, n.Accuracy, string(n.Result))
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
