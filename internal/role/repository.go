package role

import (
	"context"

	"github.com/google/uuid"
)

type RoleRepository interface {
	Create(ctx context.Context, role *Role) error

	GetByID(ctx context.Context, id uuid.UUID) (*Role, error)

	ListByApplication(ctx context.Context, applicationID uuid.UUID) ([]*Role, error)

	Update(ctx context.Context, role *Role) error

	Delete(ctx context.Context, id uuid.UUID) error
}
