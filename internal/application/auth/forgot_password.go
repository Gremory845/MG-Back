package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"mg-back/internal/domain/user"
	"mg-back/internal/infrastructure/security"
)

type ForgotPasswordUseCase struct {
	users  user.UserRepository
	resets user.PasswordResetRepository
	mailer Mailer
	secret []byte
}

func NewForgotPasswordUseCase(users user.UserRepository, resets user.PasswordResetRepository, mailer Mailer, secret []byte) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{users: users, resets: resets, mailer: mailer, secret: secret}
}

// Execute SIEMPRE devuelve nil cuando el correo no existe o se aplicó un límite:
// la respuesta es idéntica exista o no la cuenta, para que nadie pueda usar
// este endpoint para averiguar qué correos están registrados.
func (uc *ForgotPasswordUseCase) Execute(ctx context.Context, input ForgotPasswordInput) error {
	u, err := uc.users.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil
		}
		return err
	}

	now := time.Now().UTC()

	// Límites contra abuso (spam de correos)
	recent, err := uc.resets.CountCreatedSince(ctx, u.ID, now.Add(-resetCooldown))
	if err != nil {
		return err
	}
	if recent > 0 {
		return nil
	}
	lastHour, err := uc.resets.CountCreatedSince(ctx, u.ID, now.Add(-time.Hour))
	if err != nil {
		return err
	}
	if lastHour >= resetMaxPerHour {
		return nil
	}

	code, err := security.GenerateResetCode()
	if err != nil {
		return err
	}

	// Un código nuevo invalida los anteriores
	if err := uc.resets.InvalidateAllForUser(ctx, u.ID); err != nil {
		return err
	}
	err = uc.resets.Create(ctx, &user.PasswordResetCode{
		ID:        uuid.NewString(),
		UserID:    u.ID,
		CodeHash:  security.HashResetCode(uc.secret, string(u.ID), code),
		ExpiresAt: now.Add(resetCodeTTL),
		CreatedAt: now,
	})
	if err != nil {
		return err
	}

	// El correo se envía en segundo plano: así la respuesta tarda lo mismo
	// exista o no la cuenta (enviar por SMTP tarda cientos de ms y delataría
	// a las cuentas reales por tiempo de respuesta).
	email, username := u.Email, u.Username
	go func() {
		sendCtx, cancel := context.WithTimeout(context.Background(), resetMailTimeout)
		defer cancel()
		if err := uc.mailer.SendPasswordResetCode(sendCtx, email, username, code, int(resetCodeTTL.Minutes())); err != nil {
			log.Printf("ERROR enviando el correo de recuperación (usuario %s): %v", u.ID, err)
		}
	}()

	return nil
}
