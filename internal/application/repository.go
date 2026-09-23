package application

import (
	"context"

	"github.com/google/uuid"
)

type ApplicationRepository interface {
	Create(ctx context.Context,application *Application)error

	GetByID(ctx context.Context,id uuid.UUID)(*Application,error)

	GetBySlug(ctx context.Context,slug string)(*Application,error)

	List(ctx context.Context)([]*Application,error)
}