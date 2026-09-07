package token

import "time"

type TokenRepository interface {
	Create(token *Token) error
	GetByID(id string) (*Token, error)
	Revoke(id string) error
	ConsumeRefreshToken(id string, now time.Time) (*Token, error)
}
