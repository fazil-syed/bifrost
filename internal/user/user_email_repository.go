package user

import (
	"context"

	"github.com/google/uuid"
)

type UserEmailRepository interface {
	Create(ctx context.Context, userEmail *UserEmail) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*UserEmail, error)
	GetByEmail(ctx context.Context, email Email) (*UserEmail, error)
}
