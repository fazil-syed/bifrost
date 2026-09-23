package application

import (
	"context"

	"github.com/google/uuid"
)

type ApplicationService interface {
	Create(ctx context.Context, name string, slug string, ownerUserID *uuid.UUID, ownerTeamID *uuid.UUID) (*Application, error)

	GetByID(ctx context.Context, id uuid.UUID) (*Application, error)

	GetBySlug(ctx context.Context, slug string) (*Application, error)

	List(ctx context.Context) ([]*Application, error)
}
