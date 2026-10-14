package game

import (
	"context"
	"time"

	"github.com/google/uuid"

	"mg-back/internal/domain/game"
	"mg-back/internal/domain/music"
)

const pointsPerHit = 100

type SubmitAttemptUseCase struct {
	levels   music.LevelRepository
	attempts game.GameAttemptRepository
}

func NewSubmitAttemptUseCase(levels music.LevelRepository, attempts game.GameAttemptRepository) *SubmitAttemptUseCase {
	return &SubmitAttemptUseCase{levels: levels, attempts: attempts}
}

func (uc *SubmitAttemptUseCase) Execute(ctx context.Context, userID string, input SubmitAttemptInput) (*SubmitAttemptOutput, error) {
	level, err := uc.levels.FindByID(ctx, input.LevelID)
	if err != nil {
		return nil, err
	}

	levelNotes, err := uc.levels.FindNotesByLevelID(ctx, input.LevelID)
	if err != nil {
		return nil, err
	}

	// Validación: el cliente debe reportar exactamente las notas del nivel
	if len(input.Notes) != len(levelNotes) {
		return nil, game.ErrGameAttemptInvalid
	}

	notesByPosition := make(map[int]music.LevelNote, len(levelNotes))
	for _, ln := range levelNotes {
		notesByPosition[ln.Position] = ln
	}

	hits := 0
	noteAttempts := make([]game.NoteAttempt, 0, len(input.Notes))

	for _, ni := range input.Notes {
		ln, ok := notesByPosition[ni.Position]
		if !ok {
			return nil, game.ErrGameAttemptInvalid
		}

		result := game.NoteResult(ni.Result)
		noteAccuracy := 0.0
		if result == game.ResultHit {
			hits++
			noteAccuracy = 100.0
		}

		// La frecuencia objetivo la calcula el servidor desde el nivel,
		// no la confiamos al cliente.
		targetFrequency := ln.Pitch.Frequency()

		var cents *float64
		if ni.DetectedFrequency != nil {
			c := music.CentsDifference(targetFrequency, *ni.DetectedFrequency)
			cents = &c
		}

		noteAttempts = append(noteAttempts, game.NoteAttempt{
			Position:           ni.Position,
			TargetPitchClass:   ln.Pitch.PitchClass,
			TargetOctave:       ln.Pitch.Octave,
			DetectedPitchClass: ni.DetectedPitchClass,
			DetectedOctave:     ni.DetectedOctave,
			TargetFrequency:    targetFrequency,
			DetectedFrequency:  ni.DetectedFrequency,
			CentsDifference:    cents,
			Accuracy:           noteAccuracy,
			Result:             result,
		})
	}

	total := len(input.Notes)
	mistakes := total - hits
	accuracy := (float64(hits) / float64(total)) * 100
	score := hits * pointsPerHit

	// Si el cliente dice que completó el nivel, la precisión debe alcanzar
	// el mínimo configurado para ese nivel.
	if input.Completed && accuracy < level.RequiredAccuracy {
		return nil, game.ErrGameAttemptInvalid
	}

	now := time.Now().UTC()
	var finishedAt *time.Time
	if input.Completed {
		finishedAt = &now
	}

	attempt := &game.GameAttempt{
		ID:         uuid.NewString(),
		UserID:     userID,
		LevelID:    input.LevelID,
		Score:      score,
		Accuracy:   accuracy,
		Mistakes:   mistakes,
		Completed:  input.Completed,
		DurationMs: input.DurationMs,
		StartedAt:  now,
		FinishedAt: finishedAt,
		Notes:      noteAttempts,
	}

	if err := uc.attempts.Submit(ctx, attempt); err != nil {
		return nil, err
	}

	return &SubmitAttemptOutput{
		ID:        attempt.ID,
		Score:     attempt.Score,
		Accuracy:  attempt.Accuracy,
		Mistakes:  attempt.Mistakes,
		Completed: attempt.Completed,
	}, nil
}
