package token

import (
	"time"

	"github.com/google/uuid"
)

type RefreshTokenFamily struct {
	ID        uuid.UUID
	CreatedAt time.Time
	RevokedAt *time.Time
}
