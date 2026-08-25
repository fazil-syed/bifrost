package user

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Enable(ctx context.Context, id uuid.UUID) error
	Disable(ctx context.Context, id uuid.UUID) error
}
