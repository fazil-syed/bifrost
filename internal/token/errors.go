package token

import "errors"

var (
	ErrTokenNotFound      = errors.New("token not found")
	ErrTokenRevoked       = errors.New("token revoked")
	ErrTokenExpired       = errors.New("token expired")
	ErrRefreshTokenUsed   = errors.New("refresh token already used")
	ErrRefreshTokenReplay = errors.New("refresh token replay detected")
	ErrInvalidTokenType   = errors.New("invalid token type")

	ErrTokenRepositoryRequired        = errors.New("token repository is required")
	ErrAuthorizationEvaluatorRequired = errors.New("authorization evaluator is required")
	ErrInvalidAccessTokenLifetime     = errors.New("access token lifetime must be greater than zero")
	ErrInvalidRefeshTokenLifetime     = errors.New("refresh token lifetime must be greater than zero")
	ErrInvalidAudience                = errors.New("invalid token audience")
	ErrInsufficientScope              = errors.New("insufficient scope")
)
