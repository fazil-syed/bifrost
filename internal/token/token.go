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

type RefreshTokenStatus string

const (
	RefreshTokenActive  RefreshTokenStatus = "ACTIVE"
	RefreshTokenUsed    RefreshTokenStatus = "USED"
	RefreshTokenRevoked RefreshTokenStatus = "REVOKED"
)

type Token struct {
	ID            string
	Type          Type
	UserID        uuid.UUID
	TenantID      uuid.UUID
	ApplicationID uuid.UUID
	Scopes        []string
	Audience      string
	IssuedAt      time.Time
	ExpiresAt     time.Time
	RevokedAt     *time.Time

	// Refresh-token rotation state
	//These fields are only valid for refresh tokens
	FamilyID uuid.UUID
	Status   RefreshTokenStatus
	UsedAt   *time.Time
}

func NewAccessToken(tenantID uuid.UUID, applicationID uuid.UUID, audience string, userID uuid.UUID, scopes []string, now time.Time, lifetime time.Duration) (*Token, error) {
	token, err := newToken(TypeAccess, userID, tenantID, applicationID, audience, scopes, now, lifetime)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func NewRefreshToken(tenantID uuid.UUID, applicationID uuid.UUID, audience string, userID uuid.UUID, scopes []string, familyID uuid.UUID, now time.Time, lifetime time.Duration) (*Token, error) {
	if familyID == uuid.Nil {
		familyID = uuid.New()
	}
	token, err := newToken(TypeRefresh, userID, tenantID, applicationID, audience, scopes, now, lifetime)
	if err != nil {
		return nil, err
	}
	token.FamilyID = familyID
	token.Status = RefreshTokenActive

	return token, nil
}

func newToken(tokenType Type, userID uuid.UUID, tenantID uuid.UUID, applicationID uuid.UUID, audience string, scopes []string, now time.Time, lifetime time.Duration) (*Token, error) {
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
	if audience == "" {
		return nil, fmt.Errorf("audience is required")
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
		Scopes:        cloneScoppes(scopes),
		Audience:      audience,
		IssuedAt:      now,
		ExpiresAt:     now.Add(lifetime),
	}, nil
}

func cloneScoppes(scopes []string) []string {
	if len(scopes) == 0 {
		return nil
	}
	cloned := make([]string, len(scopes))
	copy(cloned, scopes)
	return cloned
}
