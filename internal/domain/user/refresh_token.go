package user

import "time"

// RefreshToken representa una sesión de larga duración. El dominio solo
// maneja el hash del token, nunca el valor en texto plano — ese se
// genera y se entrega al cliente en la capa de aplicación/seguridad,
// y nunca se persiste.
type RefreshToken struct {
	ID        string
	UserID    ID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
