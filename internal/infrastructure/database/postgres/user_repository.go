package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"mg-back/internal/domain/user"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) FindByID(ctx context.Context, id user.ID) (*user.User, error) {
	return r.findOne(ctx,
		`SELECT id, email, username, password_hash, is_active, created_at, updated_at
		 FROM users WHERE id = $1`, id)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return r.findOne(ctx,
		`SELECT id, email, username, password_hash, is_active, created_at, updated_at
		 FROM users WHERE email = $1`, email)
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	return r.findOne(ctx,
		`SELECT id, email, username, password_hash, is_active, created_at, updated_at
		 FROM users WHERE username = $1`, username)
}

func (r *UserRepository) findOne(ctx context.Context, query string, arg interface{}) (*user.User, error) {
	row := r.pool.QueryRow(ctx, query, arg)

	var u user.User
	var id string
	err := row.Scan(&id, &u.Email, &u.Username, &u.PasswordHash, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	u.ID = user.ID(id)
	return &u, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id user.ID, passwordHash string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1
	`, id, passwordHash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return user.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (id, email, username, password_hash, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, u.ID, u.Email, u.Username, u.PasswordHash, u.IsActive, u.CreatedAt, u.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			switch pgErr.ConstraintName {
			case "users_email_unique":
				return user.ErrEmailAlreadyExists
			case "users_username_unique":
				return user.ErrUsernameTaken
			}
		}
		return err
	}

	return nil
}
