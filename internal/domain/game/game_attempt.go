package game

import "time"

type NoteResult string

const (
	ResultHit  NoteResult = "HIT"
	ResultMiss NoteResult = "MISS"
)

type NoteAttempt struct {
	Position           int
	TargetPitchClass   int
	TargetOctave       int
	DetectedPitchClass *int
	DetectedOctave     *int
	TargetFrequency    float64
	DetectedFrequency  *float64
	CentsDifference    *float64
	Accuracy           float64
	Result             NoteResult
}

type GameAttempt struct {
	ID         string
	UserID     string
	LevelID    string
	Score      int
	Accuracy   float64
	Mistakes   int
	Completed  bool
	DurationMs *int
	StartedAt  time.Time
	FinishedAt *time.Time
	Notes      []NoteAttempt
}
