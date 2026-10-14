package scale

import (
	"context"

	"mg-back/internal/domain/music"
)

type GetLevelUseCase struct {
	levels music.LevelRepository
}

func NewGetLevelUseCase(levels music.LevelRepository) *GetLevelUseCase {
	return &GetLevelUseCase{levels: levels}
}

func (uc *GetLevelUseCase) Execute(ctx context.Context, levelID string) (*LevelDetailOutput, error) {
	level, err := uc.levels.FindByID(ctx, levelID)
	if err != nil {
		return nil, err
	}

	notes, err := uc.levels.FindNotesByLevelID(ctx, levelID)
	if err != nil {
		return nil, err
	}

	noteOutputs := make([]LevelNoteOutput, 0, len(notes))
	for _, n := range notes {
		noteOutputs = append(noteOutputs, LevelNoteOutput{
			Position:       n.Position,
			PitchClass:     n.Pitch.PitchClass,
			Octave:         n.Pitch.Octave,
			NoteName:       n.Pitch.NoteName(),
			Frequency:      n.Pitch.Frequency(),
			ToleranceCents: n.ToleranceCents,
		})
	}

	return &LevelDetailOutput{
		ID:               level.ID,
		ScaleID:          level.ScaleID,
		LevelNumber:      level.LevelNumber,
		Name:             level.Name,
		RequiredAccuracy: level.RequiredAccuracy,
		Notes:            noteOutputs,
	}, nil
}
