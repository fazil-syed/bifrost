package authentication

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/fazil-syed/bifrost/internal/user"
)

func NewService(users user.UserService) (AuthenticationService, error) {
	dummyPassword, err := generateDummyPassword()
	if err != nil {
		return nil, fmt.Errorf("generate dummy password: %w", err)
	}

	dummyHash, err := user.HashPassword(dummyPassword)

	if err != nil {
		return nil, fmt.Errorf("create dummy password hash: %w", err)
	}

	return &service{
		users:     users,
		dummyHash: dummyHash,
	}, nil
}

func (s *service) AuthenticatePassword(ctx context.Context, email string, password string) (*Principal, error) {
	bifrostUser, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			_, _ = user.VerifyPassword(password, s.dummyHash)
			return nil, ErrAuthenticationFailed
		}
		return nil, fmt.Errorf("authenticate user: %w", err)
	}
	if bifrostUser.Status == user.UserStatusDisabled {
		_, _ = user.VerifyPassword(password, s.dummyHash)
		return nil, ErrAuthenticationFailed
	}
	valid, err := s.users.VerifyPassword(ctx, bifrostUser.ID, password)

	if err != nil {
		return nil, fmt.Errorf("verify user password: %w", err)
	}

	if !valid {
		return nil, ErrAuthenticationFailed
	}
	return &Principal{
		UserID:               bifrostUser.ID,
		AuthenticationMethod: AuthenticationMethodPassword,
		AuthenticatedAt:      time.Now(),
	}, nil

}

func generateDummyPassword() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}
