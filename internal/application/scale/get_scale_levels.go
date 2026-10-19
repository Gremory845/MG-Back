package scale

import (
	"context"

	"mg-back/internal/domain/music"
)

type LevelSummary struct {
	ID               string  `json:"id"`
	LevelNumber      int     `json:"level_number"`
	Name             string  `json:"name"`
	RequiredAccuracy float64 `json:"required_accuracy"`
}

type ScaleLevelsOutput struct {
	ScaleID   string         `json:"scale_id"`
	ScaleName string         `json:"scale_name"`
	Category  string         `json:"category"`
	Family    string         `json:"family"`
	Levels    []LevelSummary `json:"levels"`
}

type GetScaleLevelsUseCase struct {
	scales music.ScaleRepository
	levels music.LevelRepository
}

func NewGetScaleLevelsUseCase(scales music.ScaleRepository, levels music.LevelRepository) *GetScaleLevelsUseCase {
	return &GetScaleLevelsUseCase{scales: scales, levels: levels}
}

func (uc *GetScaleLevelsUseCase) Execute(ctx context.Context, scaleID string) (*ScaleLevelsOutput, error) {
	sc, err := uc.scales.FindByID(ctx, scaleID)
	if err != nil {
		return nil, err
	}

	levels, err := uc.levels.FindByScaleID(ctx, scaleID)
	if err != nil {
		return nil, err
	}

	out := &ScaleLevelsOutput{ScaleID: sc.ID, ScaleName: sc.Name, Category: sc.Category, Family: sc.Family, Levels: make([]LevelSummary, 0, len(levels))}
	for _, l := range levels {
		out.Levels = append(out.Levels, LevelSummary{
			ID:               l.ID,
			LevelNumber:      l.LevelNumber,
			Name:             l.Name,
			RequiredAccuracy: l.RequiredAccuracy,
		})
	}
	return out, nil
}
