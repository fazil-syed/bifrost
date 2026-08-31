package user

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, email Email) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)

	GetByExternalIdentity(ctx context.Context, issuer string, subject string) (*User, error)
	LinkExternalIdentity(ctx context.Context, userID uuid.UUID, issuer string, subject string) error
	Enable(ctx context.Context, id uuid.UUID) error
	Disable(ctx context.Context, id uuid.UUID) error

	SetPassword(ctx context.Context, userID uuid.UUID, password string) error
	ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword string, newPassword string) error
	VerifyPassword(ctx context.Context, userID uuid.UUID, password string) (bool, error)
}
