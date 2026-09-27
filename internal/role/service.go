package role

import (
	"context"

	"github.com/google/uuid"
)

type RoleService interface {
	Create(ctx context.Context, applicationID uuid.UUID, name string) (*Role, error)

	GetByID(ctx context.Context, id uuid.UUID) (*Role, error)

	ListByApplication(ctx context.Context, applicationID uuid.UUID) ([]*Role, error)

	Update(ctx context.Context, id uuid.UUID, name string) error

	Delete(ctx context.Context, id uuid.UUID) error
}
