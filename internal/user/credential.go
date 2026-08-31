package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type PasswordCredential struct {
	UserID       uuid.UUID
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var ErrPasswordCredentialNotFound = errors.New("password credential not found")
var ErrInvalidPassword = errors.New("invalid password")

func NewPasswordCredential(userID uuid.UUID, passwordHash string, now time.Time) *PasswordCredential {
	return &PasswordCredential{
		UserID:       userID,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
