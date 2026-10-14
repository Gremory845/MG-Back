package game

type NoteAttemptInput struct {
	Position           int      `json:"position" validate:"required,min=1"`
	DetectedPitchClass *int     `json:"detected_pitch_class" validate:"omitempty,min=0,max=11"`
	DetectedOctave     *int     `json:"detected_octave" validate:"omitempty,min=0,max=10"`
	DetectedFrequency  *float64 `json:"detected_frequency" validate:"omitempty,gt=0"`
	Result             string   `json:"result" validate:"required,oneof=HIT MISS"`
}

type SubmitAttemptInput struct {
	LevelID    string             `json:"level_id" validate:"required,uuid"`
	Completed  bool               `json:"completed"`
	DurationMs *int               `json:"duration_ms" validate:"omitempty,min=0"`
	Notes      []NoteAttemptInput `json:"notes" validate:"required,min=1,dive"`
}

type SubmitAttemptOutput struct {
	ID        string  `json:"id"`
	Score     int     `json:"score"`
	Accuracy  float64 `json:"accuracy"`
	Mistakes  int     `json:"mistakes"`
	Completed bool    `json:"completed"`
}
