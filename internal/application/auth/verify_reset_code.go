package auth

import (
	"context"
	"errors"
	"time"

	"mg-back/internal/domain/user"
	"mg-back/internal/infrastructure/security"
)

type VerifyResetCodeUseCase struct {
	users  user.UserRepository
	resets user.PasswordResetRepository
	secret []byte
}

func NewVerifyResetCodeUseCase(users user.UserRepository, resets user.PasswordResetRepository, secret []byte) *VerifyResetCodeUseCase {
	return &VerifyResetCodeUseCase{users: users, resets: resets, secret: secret}
}

// Execute comprueba el código de 6 dígitos y, si es correcto, entrega un token
// de un solo uso para poder escribir la nueva contraseña. Todos los fallos
// (correo inexistente, código incorrecto, vencido o con demasiados intentos)
// devuelven el MISMO error.
func (uc *VerifyResetCodeUseCase) Execute(ctx context.Context, input VerifyResetCodeInput) (*VerifyResetCodeOutput, error) {
	u, err := uc.users.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, user.ErrInvalidResetCode
		}
		return nil, err
	}

	pending, err := uc.resets.FindLatestPending(ctx, u.ID)
	if err != nil {
		return nil, err // ErrInvalidResetCode si no hay solicitud
	}

	now := time.Now().UTC()
	if now.After(pending.ExpiresAt) || pending.Attempts >= resetMaxAttempts {
		return nil, user.ErrInvalidResetCode
	}

	// Contamos el intento ANTES de comparar y de forma atómica: así ni siquiera
	// con muchas peticiones en paralelo se pueden probar más de 5 códigos.
	attempts, err := uc.resets.IncrementAttempts(ctx, pending.ID)
	if err != nil {
		return nil, err
	}
	if attempts > resetMaxAttempts {
		return nil, user.ErrInvalidResetCode
	}

	if !security.ResetCodeMatches(uc.secret, string(u.ID), input.Code, pending.CodeHash) {
		return nil, user.ErrInvalidResetCode
	}

	plain, hash, err := security.GenerateResetToken()
	if err != nil {
		return nil, err
	}
	if err := uc.resets.MarkVerified(ctx, pending.ID, hash, now.Add(resetTokenTTL)); err != nil {
		return nil, err
	}

	return &VerifyResetCodeOutput{ResetToken: plain}, nil
}
