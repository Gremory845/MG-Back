package auth

import (
	"context"
	"time"

	"mg-back/internal/domain/user"
	"mg-back/internal/infrastructure/security"
)

type ResetPasswordUseCase struct {
	users         user.UserRepository
	resets        user.PasswordResetRepository
	refreshTokens user.RefreshTokenRepository
}

func NewResetPasswordUseCase(users user.UserRepository, resets user.PasswordResetRepository, refreshTokens user.RefreshTokenRepository) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{users: users, resets: resets, refreshTokens: refreshTokens}
}

func (uc *ResetPasswordUseCase) Execute(ctx context.Context, input ResetPasswordInput) error {
	rc, err := uc.resets.FindByResetTokenHash(ctx, security.HashToken(input.ResetToken))
	if err != nil {
		return err // ErrInvalidResetToken
	}

	if rc.UsedAt != nil || rc.VerifiedAt == nil || rc.ResetExpiresAt == nil || time.Now().UTC().After(*rc.ResetExpiresAt) {
		return user.ErrInvalidResetToken
	}

	// El token es de un solo uso: se reclama de forma atómica antes de hacer nada.
	claimed, err := uc.resets.Consume(ctx, rc.ID)
	if err != nil {
		return err
	}
	if !claimed {
		return user.ErrInvalidResetToken
	}

	passwordHash, err := security.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}
	if err := uc.users.UpdatePassword(ctx, rc.UserID, passwordHash); err != nil {
		return err
	}

	// Cualquier otro código pendiente deja de servir...
	if err := uc.resets.InvalidateAllForUser(ctx, rc.UserID); err != nil {
		return err
	}
	// ...y se cierran todas las sesiones abiertas: si alguien más tenía acceso a la
	// cuenta (por eso se recupera la contraseña), pierde el acceso.
	return uc.refreshTokens.RevokeAllForUser(ctx, rc.UserID)
}
