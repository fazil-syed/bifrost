package token

import (
	"time"

	"github.com/google/uuid"
)

type AuthorizationEvaluator interface {
	Evaluate(tenantID uuid.UUID, applicationID uuid.UUID, userID uuid.UUID, requestedScopes []string) ([]string, error)
}

type TokenService interface {
	Issue(
		tenantID uuid.UUID,
		applicationID uuid.UUID,
		userID uuid.UUID,
		requestedScopes []string,
		now time.Time,
	) (*TokenPair, error)

	ValidateAccessToken(
		tokenID string,
		tenantID uuid.UUID,
		applicationID uuid.UUID,
		requiredScopes []string,
		now time.Time,
	) (*Token, error)
	Refresh(
		refreshTokenID string,
		now time.Time,
	) (*TokenPair, error)
}

type TokenPair struct {
	AccessToken  *Token
	RefreshToken *Token
}
