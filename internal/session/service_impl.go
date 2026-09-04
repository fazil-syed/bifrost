package session

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type sessionService struct {
	repository SessionRepository
	lifetime   time.Duration
}

func NewSessionService(repository SessionRepository, lifetime time.Duration) (SessionService, error) {

	if repository == nil {
		return nil, fmt.Errorf("session repository is nil")
	}

	if lifetime <= 0 {
		return nil, fmt.Errorf("session lifetime must be greater than 0")
	}

	return &sessionService{
		repository: repository,
		lifetime:   lifetime,
	}, nil
}

func (s *sessionService) Create(ctx context.Context, userID uuid.UUID, authenticationMethod string, now time.Time) (*Session, error) {
	session, err := New(userID, authenticationMethod, now, s.lifetime)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	return session, nil
}

func (s *sessionService) GetByID(ctx context.Context, id string, now time.Time) (*Session, error) {
	if id == "" {
		return nil, ErrSessionNotFound
	}

	session, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if session.RevokedAt != nil {
		return nil, ErrSessionRevoked
	}

	if !now.Before(session.ExpiresAt) {

		return nil, ErrSessionExpired
	}
	return session, nil
}

func (s *sessionService) Revoke(ctx context.Context, id string, now time.Time) error {

	if id == "" {
		return ErrSessionNotFound
	}
	if err := s.repository.Revoke(ctx, id, now); err != nil {
		return err
	}
	return nil
}
