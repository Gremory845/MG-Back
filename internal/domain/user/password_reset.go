package user

import (
	"context"
	"time"
)

// PasswordResetCode representa una solicitud de recuperación de contraseña.
// El dominio solo maneja hashes: ni el código de 6 dígitos ni el token de
// restablecimiento se guardan en claro.
type PasswordResetCode struct {
	ID             string
	UserID         ID
	CodeHash       string
	ExpiresAt      time.Time
	Attempts       int
	VerifiedAt     *time.Time
	ResetTokenHash *string
	ResetExpiresAt *time.Time
	UsedAt         *time.Time
	CreatedAt      time.Time
}

type PasswordResetRepository interface {
	Create(ctx context.Context, c *PasswordResetCode) error

	// CountCreatedSince cuenta las solicitudes del usuario desde "since" (para limitar el envío de correos).
	CountCreatedSince(ctx context.Context, userID ID, since time.Time) (int, error)

	// FindLatestPending devuelve la solicitud más reciente sin verificar ni usar.
	// Si no hay ninguna devuelve ErrInvalidResetCode.
	FindLatestPending(ctx context.Context, userID ID) (*PasswordResetCode, error)

	// IncrementAttempts suma un intento de forma atómica y devuelve el total.
	IncrementAttempts(ctx context.Context, id string) (int, error)

	MarkVerified(ctx context.Context, id string, resetTokenHash string, resetExpiresAt time.Time) error

	// FindByResetTokenHash devuelve ErrInvalidResetToken si no existe.
	FindByResetTokenHash(ctx context.Context, tokenHash string) (*PasswordResetCode, error)

	// Consume marca la solicitud como usada de forma atómica.
	// Devuelve false si ya había sido usada (el token es de un solo uso).
	Consume(ctx context.Context, id string) (bool, error)

	// InvalidateAllForUser marca como usadas todas las solicitudes pendientes del usuario.
	InvalidateAllForUser(ctx context.Context, userID ID) error
}
