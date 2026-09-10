package token

import "time"

type TokenRepository interface {
	CreateToken(token *Token) error
	GetTokenByID(id string) (*Token, error)
	RevokeToken(id string) error
	CreateRefreshTokenFamily(family *RefreshTokenFamily) error
	GetRefreshTokenFamilyById(id string) (*RefreshTokenFamily, error)

	CreateInitialTokenPair(
		family *RefreshTokenFamily,
		accessToken *Token,
		refreshToken *Token,
	) error

	RotateRefreshToken(
		familyID string,
		presentedRefreshTokenID string,
		accessToken *Token,
		refreshToken *Token,
		now time.Time,
	) error
}
