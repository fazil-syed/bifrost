package tenant

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TenantRepository interface {
	Create(ctx context.Context, tenant *Tenant) error
	GetBydID(ctx context.Context, id uuid.UUID) (*Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*Tenant, error)

	List(ctx context.Context) ([]*Tenant, error)

	UpdateStatus(ctx context.Context, id uuid.UUID, status Status, updatedAt time.Time) error
}
