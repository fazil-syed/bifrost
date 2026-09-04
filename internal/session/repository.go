package session

import (
	"context"
	"time"
)

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, id string) (*Session, error)
	Revoke(ctx context.Context, id string, revokedAt time.Time) error
}
