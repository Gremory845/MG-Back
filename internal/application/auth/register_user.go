package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"mg-back/internal/domain/user"
	"mg-back/internal/infrastructure/security"
)

type RegisterUserUseCase struct {
	users user.UserRepository
}

func NewRegisterUserUseCase(users user.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{users: users}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	if _, err := uc.users.FindByEmail(ctx, input.Email); err == nil {
		return nil, user.ErrEmailAlreadyExists
	} else if !errors.Is(err, user.ErrUserNotFound) {
		return nil, err
	}

	if _, err := uc.users.FindByUsername(ctx, input.Username); err == nil {
		return nil, user.ErrUsernameTaken
	} else if !errors.Is(err, user.ErrUserNotFound) {
		return nil, err
	}

	passwordHash, err := security.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	newUser := &user.User{
		ID:           user.ID(uuid.NewString()),
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: passwordHash,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.users.Save(ctx, newUser); err != nil {
		return nil, err
	}

	return &RegisterUserOutput{
		ID:       string(newUser.ID),
		Email:    newUser.Email,
		Username: newUser.Username,
	}, nil
}
