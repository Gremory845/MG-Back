package scale

type ScaleSummary struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Difficulty int    `json:"difficulty"`
}

type LevelNoteOutput struct {
	Position       int     `json:"position"`
	PitchClass     int     `json:"pitch_class"`
	Octave         int     `json:"octave"`
	NoteName       string  `json:"note_name"`
	Frequency      float64 `json:"frequency"`
	ToleranceCents float64 `json:"tolerance_cents"`
}

type LevelDetailOutput struct {
	ID               string            `json:"id"`
	ScaleID          string            `json:"scale_id"`
	LevelNumber      int               `json:"level_number"`
	Name             string            `json:"name"`
	RequiredAccuracy float64           `json:"required_accuracy"`
	Notes            []LevelNoteOutput `json:"notes"`
}
