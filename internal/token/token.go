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
	ID            string
	Type          Type
	UserID        uuid.UUID
	TenantID      uuid.UUID
	ApplicationID uuid.UUID
	Scopes        []string
	IssuedAt      time.Time
	ExpiresAt     time.Time
	RevokedAt     *time.Time
}

func New(tokenType Type, userID uuid.UUID, tenantID uuid.UUID, applicationID uuid.UUID, scopes []string, now time.Time, lifetime time.Duration) (*Token, error) {
	if tokenType != TypeAccess && tokenType != TypeRefresh {
		return nil, fmt.Errorf("invalid token type %q", tokenType)
	}
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}
	if tenantID == uuid.Nil {
		return nil, fmt.Errorf("tenant ID is required")
	}
	if applicationID == uuid.Nil {
		return nil, fmt.Errorf("application ID is required")
	}
	if lifetime <= 0 {
		return nil, fmt.Errorf("token lifetime must be greater than zero")
	}
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	return &Token{
		ID:            base64.RawURLEncoding.EncodeToString(tokenBytes),
		Type:          tokenType,
		UserID:        userID,
		TenantID:      tenantID,
		ApplicationID: applicationID,
		Scopes:        append([]string(nil), scopes...),
		IssuedAt:      now,
		ExpiresAt:     now.Add(lifetime),
	}, nil
}
