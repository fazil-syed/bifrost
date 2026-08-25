package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type userService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) UserService {
	return &userService{
		repository: repository,
	}
}

func (s *userService) Create(ctx context.Context) (*User, error) {
	user, err := New()
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	if err := s.repository.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("persist user: %w", err)
	}
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func (s *userService) Disable(ctx context.Context, id uuid.UUID) error {
	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user for disable: %w", err)
	}

	if user.Status == UserStatusDisabled {
		return nil
	}
	if err := s.repository.Disable(ctx, user.ID, time.Now()); err != nil {
		return fmt.Errorf("disable user: %w", err)
	}
	return nil
}
func (s *userService) Enable(ctx context.Context, id uuid.UUID) error {
	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user for enable: %w", err)
	}

	if user.Status == UserStatusActive {
		return nil
	}
	if err := s.repository.Enable(ctx, user.ID, time.Now()); err != nil {
		return fmt.Errorf("enable user: %w", err)
	}
	return nil
}
