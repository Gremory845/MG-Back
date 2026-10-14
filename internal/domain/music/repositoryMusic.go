package music

import "context"

type ScaleRepository interface {
	FindAll(ctx context.Context) ([]Scale, error)
	FindByID(ctx context.Context, id string) (*Scale, error)
}

type LevelRepository interface {
	FindByID(ctx context.Context, id string) (*Level, error)
	FindByScaleID(ctx context.Context, scaleID string) ([]Level, error)
	FindNotesByLevelID(ctx context.Context, levelID string) ([]LevelNote, error)
}
