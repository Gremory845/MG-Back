package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"mg-back/internal/domain/user"
)

type PasswordResetRepository struct {
	pool *pgxpool.Pool
}

func NewPasswordResetRepository(pool *pgxpool.Pool) *PasswordResetRepository {
	return &PasswordResetRepository{pool: pool}
}

const resetColumns = `id, user_id, code_hash, expires_at, attempts, verified_at, reset_token_hash, reset_expires_at, used_at, created_at`

func scanReset(row pgx.Row) (*user.PasswordResetCode, error) {
	var c user.PasswordResetCode
	var id, userID string
	err := row.Scan(&id, &userID, &c.CodeHash, &c.ExpiresAt, &c.Attempts,
		&c.VerifiedAt, &c.ResetTokenHash, &c.ResetExpiresAt, &c.UsedAt, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	c.ID = id
	c.UserID = user.ID(userID)
	return &c, nil
}

func (r *PasswordResetRepository) Create(ctx context.Context, c *user.PasswordResetCode) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO password_reset_codes (id, user_id, code_hash, expires_at, attempts, created_at)
		VALUES ($1, $2, $3, $4, 0, $5)
	`, c.ID, c.UserID, c.CodeHash, c.ExpiresAt, c.CreatedAt)
	return err
}

func (r *PasswordResetRepository) CountCreatedSince(ctx context.Context, userID user.ID, since time.Time) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx,
		`SELECT count(*) FROM password_reset_codes WHERE user_id = $1 AND created_at >= $2`,
		userID, since).Scan(&n)
	return n, err
}

func (r *PasswordResetRepository) FindLatestPending(ctx context.Context, userID user.ID) (*user.PasswordResetCode, error) {
	c, err := scanReset(r.pool.QueryRow(ctx, `
		SELECT `+resetColumns+`
		FROM password_reset_codes
		WHERE user_id = $1 AND used_at IS NULL AND verified_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`, userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrInvalidResetCode
		}
		return nil, err
	}
	return c, nil
}

func (r *PasswordResetRepository) IncrementAttempts(ctx context.Context, id string) (int, error) {
	var attempts int
	err := r.pool.QueryRow(ctx,
		`UPDATE password_reset_codes SET attempts = attempts + 1 WHERE id = $1 RETURNING attempts`,
		id).Scan(&attempts)
	return attempts, err
}

func (r *PasswordResetRepository) MarkVerified(ctx context.Context, id string, resetTokenHash string, resetExpiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE password_reset_codes
		SET verified_at = NOW(), reset_token_hash = $2, reset_expires_at = $3
		WHERE id = $1
	`, id, resetTokenHash, resetExpiresAt)
	return err
}

func (r *PasswordResetRepository) FindByResetTokenHash(ctx context.Context, tokenHash string) (*user.PasswordResetCode, error) {
	c, err := scanReset(r.pool.QueryRow(ctx, `
		SELECT `+resetColumns+` FROM password_reset_codes WHERE reset_token_hash = $1
	`, tokenHash))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrInvalidResetToken
		}
		return nil, err
	}
	return c, nil
}

func (r *PasswordResetRepository) Consume(ctx context.Context, id string) (bool, error) {
	tag, err := r.pool.Exec(ctx,
		`UPDATE password_reset_codes SET used_at = NOW() WHERE id = $1 AND used_at IS NULL`, id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() == 1, nil
}

func (r *PasswordResetRepository) InvalidateAllForUser(ctx context.Context, userID user.ID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE password_reset_codes SET used_at = NOW() WHERE user_id = $1 AND used_at IS NULL`, userID)
	return err
}
