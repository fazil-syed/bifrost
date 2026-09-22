package team

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postgresMembershipRepository struct {
	tx pgx.Tx
}

func NewMembershipRepository(tx pgx.Tx) MembershipRepository {
	return &postgresMembershipRepository{tx: tx}
}

func (r *postgresMembershipRepository) Create(ctx context.Context, membership *Membership) error {
	const query = `
		INSERT INTO team_memberships (
			team_id,
			user_id,
			membership_type,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.tx.Exec(ctx, query, membership.TeamID, membership.UserID, membership.MembershipType, membership.CreatedAt, membership.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create team membership: %w", err)
	}
	return nil
}

func (r *postgresMembershipRepository) Get(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (*Membership, error) {
	const query = `
		SELECT
			team_id,
			user_id,
			membership_type,
			created_at,
			updated_at
		FROM team_memberships
		WHERE team_id = $1 
			AND user_id = $2
	`
	var membership Membership
	err := r.tx.QueryRow(ctx, query, teamID, userID).Scan(&membership.TeamID, &membership.UserID, &membership.MembershipType, &membership.CreatedAt, &membership.UpdatedAt)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTeamMembershipNotFound
		}
		return nil, fmt.Errorf("get team membership: %w", err)
	}
	return &membership, nil
}

func (r *postgresMembershipRepository) ListByTeam(ctx context.Context, teamID uuid.UUID) ([]*Membership, error) {
	const query = `
		SELECT 
			team_id,
			user_id,
			membership_type,
			created_at,
			updated_at
		FROM team_memberships
		WHERE team_id = $1
		ORDER BY created_at,id
	`

	rows, err := r.tx.Query(ctx, query, teamID)

	if err != nil {
		return nil, fmt.Errorf("list team memberships: %w", err)
	}
	defer rows.Close()

	memberships := make([]*Membership, 0)

	for rows.Next() {
		var membership Membership

		if err := rows.Scan(&membership.TeamID, &membership.UserID, &membership.MembershipType, &membership.CreatedAt, &membership.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan team membership: %w", err)
		}
		memberships = append(memberships, &membership)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate team memberships: %w", err)
	}
	return memberships, nil
}

func (r *postgresMembershipRepository) UpdateType(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, membershipType MembershipType, updatedAt time.Time) error {

	const query = `
		UPDATE team_memberships
		SET
			membership_type = $1,
			updated_at = $2
		WHERE team_id = $3
			AND user_id = $4
	
	`

	result, err := r.tx.Exec(ctx, query, membershipType, updatedAt, teamID, userID)
	if err != nil {
		return fmt.Errorf("update team membership type: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrTeamMembershipNotFound
	}

	return nil
}

func (r *postgresMembershipRepository) Delete(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	const query = `
		DELETE FROM team_memberships
		WHERE team_id = $1
			AND user_id = $2
	`

	result, err := r.tx.Exec(ctx, query, teamID, userID)
	if err != nil {
		return fmt.Errorf("delete team membership: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrTeamMembershipNotFound
	}

	return nil
}

func (r *postgresMembershipRepository) CountOwners(ctx context.Context, teamID uuid.UUID) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM team_memberships
		WHERE team_id = $1
		 AND membership_type = $2
	`
	var count int
	err := r.tx.QueryRow(ctx, query, teamID, MembershipOwner).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count team owners: %w", err)
	}
	return count, nil
}
