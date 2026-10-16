package auth

import (
	"errors"
	"time"
)

var (
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}
