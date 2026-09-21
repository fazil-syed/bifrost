package tenant

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type service struct {
	db *pgxpool.Pool
}

func NewTenantService(db *pgxpool.Pool) TenantService {
	return &service{db: db}
}

func (s *service) Create(ctx context.Context, name string, slug string, databaseName string) (*Tenant, error) {
	now := time.Now().UTC()

	tenant, err := New(name, slug, databaseName, now)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)

	if err != nil {
		return nil, fmt.Errorf("begin tenant creation transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewPostgresTenantRepository(tx)

	if err := repository.Create(ctx, tenant); err != nil {

		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tenant creation: %w", err)
	}
	return tenant, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, fmt.Errorf("begin tenant read transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewPostgresTenantRepository(tx)
	return repository.GetByID(ctx, id)
}

func (s *service) GetBySlug(ctx context.Context, slug string) (*Tenant, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, fmt.Errorf("begin tenant read transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewPostgresTenantRepository(tx)
	return repository.GetBySlug(ctx, slug)
}

func (s *service) List(ctx context.Context) ([]*Tenant, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})

	if err != nil {
		return nil, fmt.Errorf("begin tenant list transaction: %w", err)
	}

	defer tx.Rollback(ctx)
	repository := NewPostgresTenantRepository(tx)
	return repository.List(ctx)
}

func (s *service) Enable(ctx context.Context, id uuid.UUID) error {
	return s.updateStatus(ctx, id, StatusActive)
}
func (s *service) Disable(ctx context.Context, id uuid.UUID) error {
	return s.updateStatus(ctx, id, StatusDisabled)
}

func (s *service) updateStatus(ctx context.Context, id uuid.UUID, status Status) error {
	now := time.Now().UTC()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tenant update status transaction: %w", err)

	}

	defer tx.Rollback(ctx)

	repository := NewPostgresTenantRepository(tx)

	if err := repository.UpdateStatus(ctx, id, status, now); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tenant status update: %w", err)
	}
	return nil
}
