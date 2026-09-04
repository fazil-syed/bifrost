package authentication

import (
	"context"
	"errors"
	"time"

	"github.com/fazil-syed/bifrost/internal/session"
	"github.com/fazil-syed/bifrost/internal/user"
	"github.com/google/uuid"
)

type AuthenticationMethod string

const (
	AuthenticationMethodPassword AuthenticationMethod = "PASSWORD"
)

var ErrAuthenticationFailed = errors.New("authentication failed")

type Principal struct {
	UserID               uuid.UUID
	AuthenticationMethod AuthenticationMethod
	AuthenticatedAt      time.Time
}

type AuthenticationService interface {
	AuthenticatePassword(
		ctx context.Context,
		email string,
		password string,
	) (*Principal, error)

	LoginWithPassword(ctx context.Context, email string, paswword string) (*Principal, *session.Session, error)
}

type service struct {
	users          user.UserService
	sessionService session.SessionService
	dummyHash      string
}
