package scale

import (
	"context"

	"mg-back/internal/domain/music"
)

type GetScalesUseCase struct {
	scales music.ScaleRepository
}

func NewGetScalesUseCase(scales music.ScaleRepository) *GetScalesUseCase {
	return &GetScalesUseCase{scales: scales}
}

func (uc *GetScalesUseCase) Execute(ctx context.Context) ([]ScaleSummary, error) {
	scales, err := uc.scales.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]ScaleSummary, 0, len(scales))
	for _, s := range scales {
		out = append(out, ScaleSummary{
			ID:         s.ID,
			Name:       s.Name,
			Difficulty: s.Difficulty,
			Category:   s.Category,
			Family:     s.Family,
		})
	}
	return out, nil
}
