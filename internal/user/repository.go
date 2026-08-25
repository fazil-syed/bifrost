package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Enable(ctx context.Context, id uuid.UUID, updatedAt time.Time) error
	Disable(ctx context.Context, id uuid.UUID, updatedAt time.Time) error
}
