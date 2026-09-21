package tenant

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type postgresTenantRepository struct {
	tx pgx.Tx
}

func NewPostgresTenantRepository(tx pgx.Tx) TenantRepository {
	return &postgresTenantRepository{tx: tx}
}

func (r *postgresTenantRepository) Create(ctx context.Context, tenant *Tenant) error {
	const query = `
		INSERT INTO tenants (
			id,
			name,
			slug,
			database_name,
			status,
			created_at,
			updated_at	
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	_, err := r.tx.Exec(ctx, query, tenant.ID, tenant.Name, tenant.Slug, tenant.DatabaseName, tenant.Status, tenant.CreatedAt, tenant.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique constraint violation error code
			switch pgErr.ConstraintName {
			case "tenants_slug_key":
				return ErrTenantSlugExists
			case "tenants_database_name_key":
				return ErrTenantDatabaseNameExists
			}
		}
		return fmt.Errorf("create tenant: %w", err)

	}
	return nil
}

func (r *postgresTenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	const query = `
		SELECT
			id,
			name,
			slug,
			database_name,
			status,
			created_at,
			updated_at
		FROM tenants
		WHERE id = $1
	`

	var tenant Tenant
	err := r.tx.QueryRow(ctx, query, id).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Slug,
		&tenant.DatabaseName,
		&tenant.Status,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTenantNotFound
		}
		return nil, fmt.Errorf("get tenant %s: %w", id, err)
	}
	return &tenant, nil
}

func (r *postgresTenantRepository) GetBySlug(ctx context.Context, slug string) (*Tenant, error) {
	const query = `
		SELECT 
			id,
			name,
			slug,
			database_name,
			status,
			created_at,
			updated_at
		FROM tenants
		WHERE slug = $1
	`
	var tenant Tenant
	err := r.tx.QueryRow(ctx, query, slug).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Slug,
		&tenant.DatabaseName,
		&tenant.Status,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTenantNotFound
		}
		return nil, fmt.Errorf("get tenant by slug %s: %w", slug, err)
	}
	return &tenant, nil

}

func (r *postgresTenantRepository) List(ctx context.Context) ([]*Tenant, error) {
	const query = `
			SELECT 
			id,
			name,
			slug,
			database_name,
			status,
			created_at,
			updated_at
		FROM tenants
		ORDER BY created_at, id
	`
	rows, err := r.tx.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list tenants: %w", err)
	}

	defer rows.Close()

	tenants := make([]*Tenant, 0)

	for rows.Next() {
		var tenant Tenant
		if err := rows.Scan(&tenant.ID, &tenant.Name, &tenant.Slug, &tenant.DatabaseName, &tenant.Status, &tenant.CreatedAt, &tenant.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan tenant: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read tenants: %w", err)
	}
	return tenants, nil
}

func (r *postgresTenantRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status Status, updatedAt time.Time) error {
	const query = `
	 UPDATE tenants
	 SET
	 	status = $1,
		updated_at = $2
	WHERE id = $3
	`

	result, err := r.tx.Exec(ctx, query, status, updatedAt, id)
	if err != nil {
		return fmt.Errorf("update tenant status %s: %w", id, err)
	}
	if result.RowsAffected() == 0 {
		return ErrTenantNotFound
	}
	return nil
}
