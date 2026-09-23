package application

import (
	"context"
	"fmt"
	"time"

	"github.com/fazil-syed/bifrost/internal/team"
	"github.com/fazil-syed/bifrost/internal/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type applicationService struct {
	db          *pgxpool.Pool
	userService user.UserService
	teamService team.TeamService
}

func NewApplicationService(db *pgxpool.Pool, userService user.UserService, teamService team.TeamService) ApplicationService {
	return &applicationService{db: db, userService: userService, teamService: teamService}
}

func (s *applicationService) Create(ctx context.Context, name, slug string, ownerUserID *uuid.UUID, ownerTeamID *uuid.UUID) (*Application, error) {
	if ownerUserID != nil {
		if _, err := s.userService.GetByID(ctx, *ownerUserID); err != nil {
			return nil, fmt.Errorf("validate application owner user: %w", err)
		}
	}

	if ownerTeamID != nil {
		if _, err := s.teamService.GetByID(ctx, *ownerTeamID); err != nil {
			return nil, fmt.Errorf("validate application owner team: %w", err)
		}
	}

	now := time.Now().UTC()

	application, err := New(name, slug, ownerUserID, ownerTeamID, now)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin create application transaction: %w", err)
	}

	defer tx.Rollback(ctx)
	repository := NewApplicationRepository(tx)

	if err := repository.Create(ctx, application); err != nil {
		return nil, fmt.Errorf("create application: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit create application transaction: %w", err)
	}

	return application, nil
}

func (s *applicationService) GetByID(ctx context.Context, id uuid.UUID) (*Application, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, fmt.Errorf("begin get application transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewApplicationRepository(tx)
	return repository.GetByID(ctx, id)
}
func (s *applicationService) GetBySlug(ctx context.Context, slug string) (*Application, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, fmt.Errorf("begin get application transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewApplicationRepository(tx)
	return repository.GetBySlug(ctx, slug)
}
func (s *applicationService) List(ctx context.Context) ([]*Application, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, fmt.Errorf("begin list application transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewApplicationRepository(tx)
	return repository.List(ctx)
}
