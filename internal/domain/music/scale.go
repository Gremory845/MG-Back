package music

// ScaleDefinition representa la fórmula de intervalos de un tipo de
// escala (ej. Major = [0,2,4,5,7,9,11,12]), independiente de la tónica.
type ScaleDefinition struct {
	ID        string
	Code      string // "MAJOR", "DORIAN", etc.
	Name      string
	Intervals []int
}

// Scale es una escala concreta: una ScaleDefinition aplicada a una
// tónica específica (ej. C Major, D Dorian).
type Scale struct {
	ID                string
	ScaleDefinitionID string
	RootPitchClass    int // 0-11
	Name              string
	Difficulty        int
	Category          string // "MAJOR" | "MINOR"
	Family            string // "NATURAL" | "HARMONIC" | "MELODIC"
}

// Level es un nivel dentro de una escala (ej. C Major - Level 1).
type Level struct {
	ID               string
	ScaleID          string
	LevelNumber      int
	Name             string
	RequiredAccuracy float64
}

// LevelNote es una nota que el jugador debe cantar dentro de un nivel.
type LevelNote struct {
	ID             string
	LevelID        string
	Position       int
	Pitch          Pitch
	ToleranceCents float64
}
