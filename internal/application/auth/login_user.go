package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"mg-back/internal/domain/user"
	"mg-back/internal/infrastructure/security"
)

type LoginUseCase struct {
	users         user.UserRepository
	refreshTokens user.RefreshTokenRepository
	jwtSecret     []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewLoginUseCase(
	users user.UserRepository,
	refreshTokens user.RefreshTokenRepository,
	jwtSecret []byte,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *LoginUseCase {
	return &LoginUseCase{
		users:         users,
		refreshTokens: refreshTokens,
		jwtSecret:     jwtSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	u, err := uc.users.FindByUsername(ctx, input.Username)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			// Mismo error que una contraseña incorrecta: no revelamos
			// si el usuario existe o no.
			return nil, user.ErrInvalidCredentials
		}
		return nil, err
	}

	valid, err := security.VerifyPassword(input.Password, u.PasswordHash)
	if err != nil || !valid {
		return nil, user.ErrInvalidCredentials
	}

	return uc.issueTokens(ctx, u.ID)
}

// issueTokens genera un access token (JWT) + un refresh token (opaco) y
// guarda el hash del refresh token. Lo comparte LoginUseCase y
// RefreshTokenUseCase para no duplicar la lógica de emisión.
func (uc *LoginUseCase) issueTokens(ctx context.Context, userID user.ID) (*LoginOutput, error) {
	accessToken, err := security.GenerateAccessToken(string(userID), uc.jwtSecret, uc.accessTTL)
	if err != nil {
		return nil, err
	}

	refreshPlain, refreshHash, err := security.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	rt := &user.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    userID,
		TokenHash: refreshHash,
		ExpiresAt: now.Add(uc.refreshTTL),
		CreatedAt: now,
	}
	if err := uc.refreshTokens.Save(ctx, rt); err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshPlain,
	}, nil
}
