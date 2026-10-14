package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"mg-back/internal/domain/music"
)

type ScaleRepository struct {
	pool *pgxpool.Pool
}

func NewScaleRepository(pool *pgxpool.Pool) *ScaleRepository {
	return &ScaleRepository{pool: pool}
}

func (r *ScaleRepository) FindAll(ctx context.Context) ([]music.Scale, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, scale_definition_id, root_pitch_class, name, difficulty
		FROM scales
		WHERE is_active = TRUE
		ORDER BY difficulty, name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scales []music.Scale
	for rows.Next() {
		var s music.Scale
		if err := rows.Scan(&s.ID, &s.ScaleDefinitionID, &s.RootPitchClass, &s.Name, &s.Difficulty); err != nil {
			return nil, err
		}
		scales = append(scales, s)
	}
	return scales, rows.Err()
}

func (r *ScaleRepository) FindByID(ctx context.Context, id string) (*music.Scale, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, scale_definition_id, root_pitch_class, name, difficulty
		FROM scales WHERE id = $1
	`, id)

	var s music.Scale
	if err := row.Scan(&s.ID, &s.ScaleDefinitionID, &s.RootPitchClass, &s.Name, &s.Difficulty); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, music.ErrScaleNotFound
		}
		return nil, err
	}
	return &s, nil
}
