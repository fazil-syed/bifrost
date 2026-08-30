package user

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, email Email) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Enable(ctx context.Context, id uuid.UUID) error
	Disable(ctx context.Context, id uuid.UUID) error
}
