package team

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type teamService struct {
	db *pgxpool.Pool
}

func NewTeamService(db *pgxpool.Pool) TeamService {
	return &teamService{db: db}
}

func (s *teamService) Create(ctx context.Context, name string, slug string, ownerUserID uuid.UUID) (*Team, error) {
	now := time.Now()
	team, err := New(name, slug, now)
	if err != nil {
		return nil, err
	}

	membership, err := NewMembership(team.ID, ownerUserID, MembershipOwner, now)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)

	if err != nil {
		return nil, fmt.Errorf("begin create team transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	teamRepository := NewTeamRepository(tx)
	membershipRepository := NewMembershipRepository(tx)

	if err := teamRepository.Create(ctx, team); err != nil {
		return nil, fmt.Errorf("create team: %w", err)
	}
	if err := membershipRepository.Create(ctx, membership); err != nil {
		return nil, fmt.Errorf("create initial team owner: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit create team transaction: %w", err)
	}
	return team, nil
}

func (s *teamService) GetByID(ctx context.Context, id uuid.UUID) (*Team, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin get team transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	teamRepository := NewTeamRepository(tx)

	return teamRepository.GetByID(ctx, id)
}
func (s *teamService) GetBySlug(ctx context.Context, slug string) (*Team, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin get team transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	teamRepository := NewTeamRepository(tx)

	return teamRepository.GetBySlug(ctx, slug)
}

func (s *teamService) List(ctx context.Context) ([]*Team, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin list teams transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	teamRepository := NewTeamRepository(tx)

	return teamRepository.List(ctx)
}
