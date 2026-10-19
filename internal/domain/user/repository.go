package user

import "context"

// UserRepository define el contrato que debe cumplir cualquier
// implementación de persistencia (Postgres, memoria para tests, etc).
// El dominio no sabe nada de SQL ni de pgx — eso vive en
// infrastructure/database/postgres.
type UserRepository interface {
	FindByID(ctx context.Context, id ID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	Save(ctx context.Context, u *User) error
	UpdatePassword(ctx context.Context, id ID, passwordHash string) error
}
