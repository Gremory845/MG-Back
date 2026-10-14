package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"mg-back/internal/domain/music"
)

type LevelRepository struct {
	pool *pgxpool.Pool
}

func NewLevelRepository(pool *pgxpool.Pool) *LevelRepository {
	return &LevelRepository{pool: pool}
}

func (r *LevelRepository) FindByID(ctx context.Context, id string) (*music.Level, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, scale_id, level_number, name, required_accuracy
		FROM levels WHERE id = $1
	`, id)

	var l music.Level
	var name *string
	if err := row.Scan(&l.ID, &l.ScaleID, &l.LevelNumber, &name, &l.RequiredAccuracy); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, music.ErrLevelNotFound
		}
		return nil, err
	}
	if name != nil {
		l.Name = *name
	}
	return &l, nil
}

func (r *LevelRepository) FindNotesByLevelID(ctx context.Context, levelID string) ([]music.LevelNote, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT position, pitch_class, octave, tolerance_cents
		FROM level_notes
		WHERE level_id = $1
		ORDER BY position
	`, levelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []music.LevelNote
	for rows.Next() {
		var n music.LevelNote
		var pitchClass, octave int
		if err := rows.Scan(&n.Position, &pitchClass, &octave, &n.ToleranceCents); err != nil {
			return nil, err
		}
		n.LevelID = levelID
		n.Pitch = music.Pitch{PitchClass: pitchClass, Octave: octave}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}

func (r *LevelRepository) FindByScaleID(ctx context.Context, scaleID string) ([]music.Level, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, scale_id, level_number, name, required_accuracy
		FROM levels
		WHERE scale_id = $1 AND is_active = TRUE
		ORDER BY level_number
	`, scaleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []music.Level
	for rows.Next() {
		var l music.Level
		var name *string
		if err := rows.Scan(&l.ID, &l.ScaleID, &l.LevelNumber, &name, &l.RequiredAccuracy); err != nil {
			return nil, err
		}
		if name != nil {
			l.Name = *name
		}
		levels = append(levels, l)
	}
	return levels, rows.Err()
}
