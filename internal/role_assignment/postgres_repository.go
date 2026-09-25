package roleassignment

import (
	"context"
	"errors"
	"fmt"

	"github.com/fazil-syed/bifrost/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postgresRoleAssignmentRepository struct {
	tx pgx.Tx
}

func (r *postgresRoleAssignmentRepository) Create(ctx context.Context, assignment *RoleAssignment) error {
	const query = `
		INSERT INTO role_assignments (
			id,
			role_id,
			user_id,
			team_id,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.tx.Exec(ctx, query, assignment.ID, assignment.RoleID, assignment.UserID, assignment.TeamID, assignment.CreatedAt, assignment.UpdatedAt)

	if err != nil {
		if _, err := database.IsUniqueViolation(err); err {
			return ErrRoleAssignmentAlreadyExists
		}
		return fmt.Errorf("create role assignment: %w", err)
	}

	return nil
}

func (r *postgresRoleAssignmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*RoleAssignment, error) {
	const query = `
		SELECT 
			id,
			role_id,
			user_id,
			team_id,
			created_at,
			updated_at
		FROM role_assignments
		WHERE id = $1
	`
	var assignment RoleAssignment
	err := r.tx.QueryRow(ctx, query, id).Scan(&assignment.ID, &assignment.RoleID, &assignment.UserID, &assignment.TeamID, &assignment.CreatedAt, &assignment.UpdatedAt)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRoleAssignmentNotFound
		}
		return nil, fmt.Errorf("get role assignment by id: %w", err)
	}
	return &assignment, nil
}

func (r *postgresRoleAssignmentRepository) GetByRoleAndUser(ctx context.Context, roleID uuid.UUID, userID uuid.UUID) (*RoleAssignment, error) {
	const query = `
		SELECT 
			id,
			role_id,
			user_id,
			team_id,
			created_at,
			updated_at
		FROM role_assignments
		WHERE role_id = $1
			AND user_id = $2
	`
	var assignment RoleAssignment
	err := r.tx.QueryRow(ctx, query, roleID, userID).Scan(&assignment.ID, &assignment.RoleID, &assignment.UserID, &assignment.TeamID, &assignment.CreatedAt, &assignment.UpdatedAt)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRoleAssignmentNotFound
		}
		return nil, fmt.Errorf("get role assignment by role and user: %w", err)
	}
	return &assignment, nil
}

func (r *postgresRoleAssignmentRepository) GetByRoleAndTeam(ctx context.Context, roleID uuid.UUID, TeamID uuid.UUID) (*RoleAssignment, error) {
	const query = `
		SELECT 
			id,
			role_id,
			user_id,
			team_id,
			created_at,
			updated_at
		FROM role_assignments
		WHERE role_id = $1
			AND team_id = $2
	`
	var assignment RoleAssignment
	err := r.tx.QueryRow(ctx, query, roleID, TeamID).Scan(&assignment.ID, &assignment.RoleID, &assignment.UserID, &assignment.TeamID, &assignment.CreatedAt, &assignment.UpdatedAt)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRoleAssignmentNotFound
		}
		return nil, fmt.Errorf("get role assignment by role and team: %w", err)
	}
	return &assignment, nil
}

func (r *postgresRoleAssignmentRepository) ListByRole(ctx context.Context, roleID uuid.UUID) ([]*RoleAssignment, error) {
	const query = `
		SELECT
			id,
			role_id,
			user_id,
			team_id,
			created_at,
			updated_at
		FROM role_assignments
		WHERE role_id = $1
	`
	return r.list(ctx, query, roleID)

}
func (r *postgresRoleAssignmentRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*RoleAssignment, error) {
	const query = `
		SELECT
			id,
			role_id,
			user_id,
			team_id,
			created_at,
			updated_at
		FROM role_assignments
		WHERE user_id = $1
	`
	return r.list(ctx, query, userID)

}
func (r *postgresRoleAssignmentRepository) ListByTeam(ctx context.Context, teamID uuid.UUID) ([]*RoleAssignment, error) {
	const query = `
		SELECT
			id,
			role_id,
			user_id,
			team_id,
			created_at,
			updated_at
		FROM role_assignments
		WHERE team_id = $1
	`
	return r.list(ctx, query, teamID)

}

func (r *postgresRoleAssignmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM role_assignments
		WHERE id = $1
	`
	result, err := r.tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete role assignment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrRoleAssignmentNotFound
	}
	return nil

}

func (r *postgresRoleAssignmentRepository) list(ctx context.Context, query string, arg uuid.UUID) ([]*RoleAssignment, error) {
	rows, err := r.tx.Query(ctx, query, arg)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	assignments := make([]*RoleAssignment, 0)

	for rows.Next() {
		var assignment RoleAssignment

		if err := rows.Scan(&assignment.ID, &assignment.RoleID, &assignment.UserID, &assignment.TeamID, &assignment.CreatedAt, &assignment.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan role assignment: %w", err)
		}

		assignments = append(assignments, &assignment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate role assignments: %w", err)
	}
	return assignments, nil

}
