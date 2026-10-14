package auth

import (
	"context"
	"time"

	"github.com/google/uuid"

	"mg-back/internal/domain/user"
	"mg-back/internal/infrastructure/security"
)

type RefreshTokenUseCase struct {
	refreshTokens user.RefreshTokenRepository
	jwtSecret     []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewRefreshTokenUseCase(
	refreshTokens user.RefreshTokenRepository,
	jwtSecret []byte,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		refreshTokens: refreshTokens,
		jwtSecret:     jwtSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, refreshTokenPlain string) (*LoginOutput, error) {
	hash := security.HashToken(refreshTokenPlain)

	stored, err := uc.refreshTokens.FindByHash(ctx, hash)
	if err != nil {
		return nil, err // ya viene como user.ErrRefreshTokenNotFound desde el repo
	}

	// Reuse detection: si este token ya estaba revocado (ya se había
	// usado antes para rotar) y alguien lo vuelve a mandar, es señal de
	// que fue robado y el atacante lo está usando después de que el
	// dueño legítimo ya rotó. Revocamos TODAS las sesiones del usuario
	// como medida defensiva.
	if stored.RevokedAt != nil {
		_ = uc.refreshTokens.RevokeAllForUser(ctx, stored.UserID)
		return nil, user.ErrRefreshTokenReused
	}

	if time.Now().UTC().After(stored.ExpiresAt) {
		return nil, user.ErrRefreshTokenExpired
	}

	// Rotación: revocamos el actual y emitimos uno nuevo
	if err := uc.refreshTokens.Revoke(ctx, stored.ID); err != nil {
		return nil, err
	}

	accessToken, err := security.GenerateAccessToken(string(stored.UserID), uc.jwtSecret, uc.accessTTL)
	if err != nil {
		return nil, err
	}

	newPlain, newHash, err := security.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	newRT := &user.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    stored.UserID,
		TokenHash: newHash,
		ExpiresAt: now.Add(uc.refreshTTL),
		CreatedAt: now,
	}
	if err := uc.refreshTokens.Save(ctx, newRT); err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: newPlain,
	}, nil
}
