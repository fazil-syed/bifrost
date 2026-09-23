package team

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type membershipService struct {
	db *pgxpool.Pool
}

func NewMembershipService(db *pgxpool.Pool) MembershipService {
	return &membershipService{db: db}
}

func (s *membershipService) AddMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	now := time.Now().UTC()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin add team member transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	teamRepository := NewTeamRepository(tx)

	membershipRepository := NewMembershipRepository(tx)

	if _, err := teamRepository.GetByID(ctx, teamID); err != nil {
		return fmt.Errorf("get team: %w", err)
	}

	membership, err := NewMembership(teamID, userID, MembershipMember, now)

	if err != nil {
		return err
	}

	if err := membershipRepository.Create(ctx, membership); err != nil {
		return fmt.Errorf("create team membership: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit add team member transaction: %w", err)
	}
	return nil
}

func (s *membershipService) GetMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (*Membership, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})

	if err != nil {
		return nil, fmt.Errorf("begin get team member transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	membershipRepository := NewMembershipRepository(tx)
	return membershipRepository.Get(ctx, teamID, userID)
}

func (s *membershipService) ListMembers(ctx context.Context, teamID uuid.UUID) ([]*Membership, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})

	if err != nil {
		return nil, fmt.Errorf("begin list team members transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	membershipRepository := NewMembershipRepository(tx)
	return membershipRepository.ListByTeam(ctx, teamID)
}

func (s *membershipService) PromoteToOwner(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	now := time.Now().UTC()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin promote team member transaction: %w", err)
	}

	defer tx.Rollback(ctx)
	teamRepository := NewTeamRepository(tx)

	membershipRepository := NewMembershipRepository(tx)

	if _, err := teamRepository.GetByID(ctx, teamID); err != nil {
		return fmt.Errorf("get team: %w", err)
	}

	membership, err := membershipRepository.Get(ctx, teamID, userID)
	if err != nil {
		return fmt.Errorf("get team membership: %w", err)
	}

	if membership.MembershipType == MembershipOwner {
		return nil
	}

	if err := membershipRepository.UpdateType(ctx, teamID, userID, MembershipOwner, now); err != nil {
		return fmt.Errorf("promote team member to owner: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit promote team member to owner transaction: %w", err)
	}

	return nil
}

func (s *membershipService) DemoteOwner(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	now := time.Now().UTC()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin demote team member transaction: %w", err)
	}

	defer tx.Rollback(ctx)
	teamRepository := NewTeamRepository(tx)

	membershipRepository := NewMembershipRepository(tx)

	if _, err := teamRepository.GetByID(ctx, teamID); err != nil {
		return fmt.Errorf("get team: %w", err)
	}

	membership, err := membershipRepository.Get(ctx, teamID, userID)
	if err != nil {
		return fmt.Errorf("get team membership: %w", err)
	}

	if membership.MembershipType != MembershipOwner {
		return nil
	}

	ownerCount, err := membershipRepository.CountOwners(ctx, teamID)
	if err != nil {
		return fmt.Errorf("count team owners: %w", err)
	}

	if ownerCount <= 1 {
		return ErrTeamMustHaveOwner
	}

	if err := membershipRepository.UpdateType(ctx, teamID, userID, MembershipMember, now); err != nil {
		return fmt.Errorf("demote team member: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit demote team member transaction: %w", err)
	}
	return nil
}

func (s *membershipService) RemoveMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin remove team member transaction: %w", err)
	}

	defer tx.Rollback(ctx)
	teamRepository := NewTeamRepository(tx)

	membershipRepository := NewMembershipRepository(tx)

	if _, err := teamRepository.GetByID(ctx, teamID); err != nil {
		return fmt.Errorf("get team: %w", err)
	}

	membership, err := membershipRepository.Get(ctx, teamID, userID)
	if err != nil {
		return fmt.Errorf("get team membership: %w", err)
	}

	if membership.MembershipType == MembershipOwner {

		ownerCount, err := membershipRepository.CountOwners(ctx, teamID)
		if err != nil {
			return fmt.Errorf("count team owners: %w", err)
		}
		if ownerCount <= 1 {
			return ErrTeamMustHaveOwner
		}

	}

	if err := membershipRepository.Delete(ctx, teamID, userID); err != nil {
		return fmt.Errorf("delete team membership: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit remove member transaction: %w", err)
	}
	return nil

}
