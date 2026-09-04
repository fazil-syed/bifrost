package session

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionService interface {
	Create(ctx context.Context, userID uuid.UUID, authenticationMethod string, now time.Time) (*Session, error)
	GetByID(ctx context.Context, id string, now time.Time) (*Session, error)
	Revoke(ctx context.Context, id string, now time.Time) error
}
