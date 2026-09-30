package token

import (
	"context"
	"fmt"
	"time"

	"github.com/fazil-syed/bifrost/internal/authorization"
	"github.com/google/uuid"
)

type tokenServiceImpl struct {
	repository           TokenRepository
	authorization        authorization.AuthorizationService
	accessTokenLifetime  time.Duration
	refreshTokenLifetime time.Duration
}

func NewTokenService(repository TokenRepository, authorization authorization.AuthorizationService, accessTokenLifetime time.Duration, refreshTokenLifetime time.Duration) (TokenService, error) {
	if repository == nil {
		return nil, ErrTokenRepositoryRequired
	}

	if authorization == nil {
		return nil, ErrAuthorizationEvaluatorRequired
	}
	if accessTokenLifetime <= 0 {
		return nil, ErrInvalidAccessTokenLifetime
	}
	if refreshTokenLifetime <= 0 {
		return nil, ErrInvalidRefeshTokenLifetime
	}

	return &tokenServiceImpl{
		repository:           repository,
		authorization:        authorization,
		accessTokenLifetime:  accessTokenLifetime,
		refreshTokenLifetime: refreshTokenLifetime,
	}, nil
}

func (s *tokenServiceImpl) Issue(ctx context.Context, tenantID uuid.UUID, applicationID uuid.UUID, userID uuid.UUID, now time.Time) (*TokenPair, error) {
	effectiveScopes, err := s.authorization.GetScopes(ctx, userID, applicationID)

	if err != nil {
		return nil, fmt.Errorf("evaluate authorization: %w", err)
	}

	audience := applicationID.String()

	accessToken, err := NewAccessToken(tenantID, applicationID, audience, userID, effectiveScopes, now, s.accessTokenLifetime)

	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	refreshToken, err := NewRefreshToken(tenantID, applicationID, audience, userID, effectiveScopes, uuid.New(), now, s.refreshTokenLifetime)

	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	family := &RefreshTokenFamily{
		ID:        refreshToken.FamilyID,
		CreatedAt: now,
	}

	if err := s.repository.CreateInitialTokenPair(
		family, accessToken, refreshToken,
	); err != nil {
		return nil, fmt.Errorf("create initial token pair: %w", err)
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *tokenServiceImpl) ValidateAccessToken(ctx context.Context, tokenID string, tenantID uuid.UUID, applicationID uuid.UUID, requiredScopes []string, now time.Time) (*Token, error) {
	token, err := s.repository.GetTokenByID(tokenID)
	if err != nil {
		return nil, err
	}
	if token.Type != TypeAccess {
		return nil, ErrInvalidTokenType
	}

	if token.TenantID != tenantID {
		return nil, ErrTokenNotFound
	}

	if token.ApplicationID != applicationID {
		return nil, ErrTokenNotFound
	}

	if token.Audience != applicationID.String() {
		return nil, ErrInvalidAudience
	}

	if token.RevokedAt != nil {
		return nil, ErrTokenRevoked
	}

	if !token.ExpiresAt.After(now) {
		return nil, ErrTokenExpired
	}

	if !containsAllScopes(token.Scopes, requiredScopes) {
		return nil, ErrInsufficientScope
	}

	return token, nil
}

func (s *tokenServiceImpl) Refresh(ctx context.Context, refreshTokenID string, now time.Time) (*TokenPair, error) {

	current, err := s.repository.GetTokenByID(refreshTokenID)

	if err != nil {
		return nil, err
	}

	if current.Type != TypeRefresh {
		return nil, ErrInvalidTokenType
	}

	if current.RevokedAt != nil {
		return nil, ErrTokenRevoked
	}

	if !current.ExpiresAt.After(now) {
		return nil, ErrTokenExpired
	}

	if current.FamilyID == uuid.Nil {
		return nil, fmt.Errorf("refresh token has no family")
	}

	effectiveScopes, err := s.authorization.GetScopes(ctx, current.UserID, current.ApplicationID)

	if err != nil {
		return nil, fmt.Errorf("evaluate authorization during refresh: %w", err)
	}

	accessToken, err := NewAccessToken(current.TenantID, current.ApplicationID, current.Audience, current.UserID, effectiveScopes, now, s.accessTokenLifetime)

	if err != nil {
		return nil, fmt.Errorf("create replacement access token: %w", err)
	}

	refreshToken, err := NewRefreshToken(current.TenantID, current.ApplicationID, current.Audience, current.UserID, effectiveScopes, current.FamilyID, now, s.refreshTokenLifetime)

	if err != nil {
		return nil, fmt.Errorf("create replacement refresh token: %w", err)
	}

	if err := s.repository.RotateRefreshToken(current.FamilyID.String(), refreshTokenID, accessToken, refreshToken, now); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *tokenServiceImpl) Revoke(ctx context.Context, tokenID string) error {
	return s.repository.RevokeToken(tokenID)
}

func containsAllScopes(granted, required []string) bool {
	if len(required) == 0 {
		return true
	}

	grantedSet := make(map[string]struct{}, len(granted))
	for _, scope := range granted {
		grantedSet[scope] = struct{}{}
	}

	for _, scope := range required {
		if _, ok := grantedSet[scope]; !ok {
			return false
		}
	}
	return true
}

func filterRequestedScopes(grantedScopes []string, requestedScopes []string) ([]string, error) {
	if len(requestedScopes) == 0 {
		return grantedScopes, nil
	}

	grantedSet := make(map[string]struct{}, len(
		grantedScopes,
	))

	for _, scope := range grantedScopes {
		grantedSet[scope] = struct{}{}
	}

	effective := make([]string, 0, len(requestedScopes))

	for _, scope := range requestedScopes {
		if _, ok := grantedSet[scope]; ok {
			effective = append(effective, scope)
		}
	}
	return effective, nil
}
