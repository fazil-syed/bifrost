package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) UserService {
	return &userService{
		db: db,
	}
}

func (s *userService) Create(ctx context.Context, email Email) (*User, error) {
	user, err := New()
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	userEmail := NewUserEmail(
		user.ID,
		email,
		user.CreatedAt,
	)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin user creation transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()
	userRpository := NewPostgresUserRepository(tx)
	userEmailRepository := NewPostgresUserEmailRepository(tx)
	if err := userRpository.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("persist user: %w", err)
	}
	if err := userEmailRepository.Create(ctx, userEmail); err != nil {
		return nil, fmt.Errorf("persist user email : %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit user creation transaction: %w", err)
	}
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin get user transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()
	userRpository := NewPostgresUserRepository(tx)
	user, err := userRpository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func (s *userService) Disable(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin user disable transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	userRpository := NewPostgresUserRepository(tx)
	user, err := userRpository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user for disable: %w", err)
	}

	if user.Status == UserStatusDisabled {
		return nil
	}
	if err := userRpository.Disable(ctx, user.ID, time.Now()); err != nil {
		return fmt.Errorf("disable user: %w", err)
	}
	return nil
}
func (s *userService) Enable(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin user enable transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()
	userRpository := NewPostgresUserRepository(tx)
	user, err := userRpository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user for enable: %w", err)
	}

	if user.Status == UserStatusActive {
		return nil
	}
	if err := userRpository.Enable(ctx, user.ID, time.Now()); err != nil {
		return fmt.Errorf("enable user: %w", err)
	}
	return nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*User, error) {
	canonicalEmail, err := NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin get user by email transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	userRepository := NewPostgresUserRepository(tx)
	user, err := userRepository.GetByEmail(ctx, canonicalEmail)
	if err != nil {
		return nil, fmt.Errorf("get user by email %s : %w", email, err)
	}
	return user, nil
}

func (s *userService) GetByExternalIdentity(ctx context.Context, issuer string, subject string) (*User, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin get user by external identity transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	externalIdentityRepository := NewPostgresExternalIdentityRepository(tx)
	externalIdentity, err := externalIdentityRepository.GetByIssuerAndSubject(ctx, issuer, subject)

	if err != nil {
		return nil, fmt.Errorf("get external identity: %w", err)
	}

	userRepository := NewPostgresUserRepository(tx)

	user, err := userRepository.GetByID(ctx, externalIdentity.UserID)

	if err != nil {
		return nil, fmt.Errorf("get user for external identity: %w", err)
	}

	return user, nil
}

func (s *userService) LinkExternalIdentity(ctx context.Context, userID uuid.UUID, issuer string, subject string) error {
	externalIdentity, err := NewExternalIdentity(userID, issuer, subject, time.Now())
	if err != nil {
		return fmt.Errorf("create external identity: %w", err)
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		fmt.Errorf("begin link external identity transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	userRepository := NewPostgresUserRepository(tx)

	if _, err := userRepository.GetByID(ctx, userID); err != nil {
		return fmt.Errorf("get user for external identity: %w", err)
	}

	externalIdentityRepository := NewPostgresExternalIdentityRepository(tx)

	if err := externalIdentityRepository.Create(ctx, *externalIdentity); err != nil {
		return fmt.Errorf("link external identity: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit link external identity transaction: %w", err)
	}

	return nil
}
