package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Type string

const (
	TypeAccess  Type = "ACCESS"
	TypeRefresh Type = "REFRESH"
)

type Token struct {
	ID        string
	Type      Type
	UserID    uuid.UUID
	SessionID string
	IssuedAt  time.Time
	ExpiresAt time.Time
	RevokedAT *time.Time
}

func New(tokenType Type, userID uuid.UUID, sessionID string, now time.Time, lifetime time.Duration) (*Token, error) {
	if tokenType != TypeAccess && tokenType != TypeRefresh {
		return nil, fmt.Errorf("invalid token type %q", tokenType)
	}
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}
	if lifetime <= 0 {
		return nil, fmt.Errorf("token lifetime must be greater than zero")
	}
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	return &Token{
		ID:        base64.RawURLEncoding.EncodeToString(tokenBytes),
		Type:      tokenType,
		UserID:    userID,
		SessionID: sessionID,
		IssuedAt:  now,
		ExpiresAt: now.Add(lifetime),
	}, nil
}
