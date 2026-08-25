package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusDisabled UserStatus = "DISABLED"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID        uuid.UUID
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New() (*User, error) {
	now := time.Now()

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        id,
		Status:    UserStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
