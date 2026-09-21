package tenant

import (
	"context"

	"github.com/google/uuid"
)

type TenantService interface {
	Create(ctx context.Context, name string, slug string, databaseName string) (*Tenant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*Tenant, error)

	List(ctx context.Context) ([]*Tenant, error)

	Enable(ctx context.Context, id uuid.UUID) error
	Disable(ctx context.Context, id uuid.UUID) error
}
