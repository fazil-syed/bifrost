package team

import (
	"context"

	"github.com/google/uuid"
)

type TeamRepository interface {
	Create(ctx context.Context, team *Team) error

	GetByID(ctx context.Context, id uuid.UUID) (*Team, error)

	GetBySlug(ctx context.Context, slug string) (*Team, error)

	List(ctx context.Context) ([]*Team, error)

	Lock(ctx context.Context, id uuid.UUID) error
}
