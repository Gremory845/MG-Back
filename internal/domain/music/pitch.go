package music

import (
	"fmt"
	"math"
)

// Pitch representa una nota musical concreta: su clase de altura
// (0=C, 1=C#, 2=D, ... 11=B) y su octava. Es un Value Object:
// inmutable y comparable por valor.
type Pitch struct {
	PitchClass int // 0-11
	Octave     int
}

var noteNames = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

// MIDINumber convierte el Pitch a su número MIDI estándar (A4 = 69, C4 = 60).
func (p Pitch) MIDINumber() int {
	return 12*(p.Octave+1) + p.PitchClass
}

// Frequency calcula la frecuencia en Hz con A4 = 440Hz y temperamento
// igual: f = 440 * 2^((n-69)/12).
func (p Pitch) Frequency() float64 {
	n := p.MIDINumber()
	return 440.0 * math.Pow(2, float64(n-69)/12.0)
}

// NoteName devuelve el nombre legible, ej. "C4".
func (p Pitch) NoteName() string {
	return fmt.Sprintf("%s%d", noteNames[p.PitchClass], p.Octave)
}

// CentsDifference calcula cuántos cents separan la frecuencia detectada
// de la objetivo. Positivo = más agudo, negativo = más grave.
func CentsDifference(targetFrequency, detectedFrequency float64) float64 {
	return 1200 * math.Log2(detectedFrequency/targetFrequency)
}
