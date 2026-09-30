package token

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TokenService interface {
	Issue(
		ctx context.Context,
		tenantID uuid.UUID,
		applicationID uuid.UUID,
		userID uuid.UUID,
		now time.Time,
	) (*TokenPair, error)

	ValidateAccessToken(
		ctx context.Context,
		tokenID string,
		tenantID uuid.UUID,
		applicationID uuid.UUID,
		requiredScopes []string,
		now time.Time,
	) (*Token, error)
	Refresh(
		ctx context.Context,
		refreshTokenID string,
		now time.Time,
	) (*TokenPair, error)

	Revoke(ctx context.Context, tokenID string) error
}

type TokenPair struct {
	AccessToken  *Token
	RefreshToken *Token
}
