package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postgresRoleRepository struct {
	tx pgx.Tx
}

func NewRoleRepository(tx pgx.Tx) RoleRepository {
	return &postgresRoleRepository{tx: tx}
}

func (r *postgresRoleRepository) Create(ctx context.Context, role *Role) error {
	const query = `
		INSERT INTO roles (
			id,
			application_id,
			name,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.tx.Exec(ctx, query, role.ID, role.ApplicationID, role.Name, role.CreatedAt, role.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create role: %w", err)
	}
	return nil
}

func (r *postgresRoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*Role, error) {
	const query = `
		SELECT 
			id,
			application_id,
			name,
			created_at,
			updated_at
		FROM roles
		WHERE id = $1
	`

	var role Role
	err := r.tx.QueryRow(ctx, query, id).Scan(&role.ID, &role.ApplicationID, &role.Name, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, fmt.Errorf("get role by id: %w", err)
	}
	return &role, nil
}

func (r *postgresRoleRepository) ListByApplication(ctx context.Context, applicationID uuid.UUID) ([]*Role, error) {
	const query = `
		SELECT
			id,
			application_id,
			name,
			created_at,
			updated_at
		FROM roles
		WHERE application_id = $1
		ORDER BY created_at,id
	`

	rows, err := r.tx.Query(ctx, query, applicationID)
	if err != nil {
		return nil, fmt.Errorf("list roles by application: %w", err)
	}

	defer rows.Close()
	roles := make([]*Role, 0)
	for rows.Next() {
		var role Role

		if err := rows.Scan(&role.ID, &role.ApplicationID, &role.Name, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan role: %w", err)
		}
		roles = append(roles, &role)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate roles: %w", err)
	}
	return roles, nil
}

func (r *postgresRoleRepository) Update(ctx context.Context, role *Role) error {
	const query = `
		UPDATE roles
		SET
			name = $1,
			updated_at = $2
		WHERE id = $3
	`

	result, err := r.tx.Exec(ctx, query, role.Name, role.UpdatedAt, role.ID)
	if err != nil {
		return fmt.Errorf("update role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrRoleNotFound
	}
	return nil
}

func (r *postgresRoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM roles 
		WHERE id = $1
	`

	result, err := r.tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrRoleNotFound
	}

	return nil
}
