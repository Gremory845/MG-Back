package user

import "time"

type ID string

type User struct {
	ID           ID
	Email        string
	Username     string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
