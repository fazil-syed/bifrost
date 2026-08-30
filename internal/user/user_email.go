package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrUserEmailNotFound = errors.New("user not found")

type UserEmail struct {
	UserID    uuid.UUID
	Email     Email
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUserEmail(userID uuid.UUID, email Email, now time.Time) *UserEmail {
	return &UserEmail{
		UserID:    userID,
		Email:     email,
		Verified:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
