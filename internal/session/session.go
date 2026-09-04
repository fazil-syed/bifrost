package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID                   string
	UserID               uuid.UUID
	AuthenticationMethod string
	CreatedAt            time.Time
	ExpiresAt            time.Time
	RevokedAt            *time.Time
}

func New(userID uuid.UUID, authenticationMethod string, now time.Time, lifetime time.Duration) (*Session, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	if authenticationMethod == "" {
		return nil, fmt.Errorf("authentication method is required")
	}

	if lifetime <= 0 {
		return nil, fmt.Errorf("session lifetime must be greater than 0")
	}

	idBytes := make([]byte, 32)
	if _, err := rand.Read(idBytes); err != nil {
		return nil, fmt.Errorf("generate session id: %w", err)
	}

	return &Session{
		ID:        base64.RawURLEncoding.EncodeToString(idBytes),
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(lifetime),
	}, nil
}
