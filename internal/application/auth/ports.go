package auth

import "context"

// Mailer es el puerto de salida para enviar correos. La implementación real
// (SMTP) vive en infrastructure/mail.
type Mailer interface {
	SendPasswordResetCode(ctx context.Context, to, username, code string, ttlMinutes int) error
}
