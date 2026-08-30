package user

import (
	"context"

	"github.com/google/uuid"
)

type PasswordCredentialRepository interface {
	Create(ctx context.Context, credential *PasswordCredential) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*PasswordCredential, error)
	Update(ctx context.Context, credential *PasswordCredential) error
}
