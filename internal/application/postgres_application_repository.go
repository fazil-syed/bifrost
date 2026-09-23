package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postgresApplicationRepository struct {
	tx pgx.Tx
}

func NewApplicationRepository(tx pgx.Tx) ApplicationRepository{
	return &postgresApplicationRepository{tx: tx}
}

func (r *postgresApplicationRepository) Create(ctx context.Context, application *Application) error {
	const query = `
		INSERT INTO applications (
			id,
			name,
			slug,
			owner_user_id,
			owner_team_id,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.tx.Exec(ctx, query, application.ID, application.Name, application.Slug, application.OwnerUserID, application.OwnerTeamID, application.CreatedAt, application.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create application: %w", err)
	}
	return nil
}

func (r *postgresApplicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*Application, error) {
	const query = `
		SELECT
			id,
			name,
			slug,
			owner_user_id,
			owner_team_id,
			created_at,
			updated_at
		FROM applications
		WHERE id = $1
	`

	var application Application

	err := r.tx.QueryRow(ctx, query, id).Scan(&application.ID, &application.Name, &application.Slug, &application.OwnerUserID, &application.OwnerTeamID, &application.CreatedAt, &application.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("get application by id: %w", err)
	}
	return &application, nil
}

func (r *postgresApplicationRepository) GetBySlug(ctx context.Context, slug string) (*Application, error) {
	const query = `
		SELECT
			id,
			name,
			slug,
			owner_user_id,
			owner_team_id,
			created_at,
			updated_at
		FROM applications
		WHERE slug = $1
	`

	var application Application

	err := r.tx.QueryRow(ctx, query, slug).Scan(&application.ID, &application.Name, &application.Slug, &application.OwnerUserID, &application.OwnerTeamID, &application.CreatedAt, &application.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("get application by slug: %w", err)
	}
	return &application, nil
}

func (r *postgresApplicationRepository) List(ctx context.Context) ([]*Application, error) {
	const query = `
		SELECT
			id,
			name,
			slug,
			owner_user_id,
			owner_team_id,
			created_at,
			updated_at
		FROM applications
		ORDER BY created_at,id
	`
	rows, err := r.tx.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list applications: %w", err)
	}

	defer rows.Close()

	applications := make([]*Application, 0)

	for rows.Next() {
		var application Application
		if err := rows.Scan(&application.ID, &application.Name, &application.Slug, &application.OwnerUserID, &application.OwnerTeamID, &application.CreatedAt, &application.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan application: %w", err)
		}
		applications = append(applications, &application)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applications: %w",err)
	}

	return applications, nil
}
