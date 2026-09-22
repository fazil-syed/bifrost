package team

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postgresTeamRepository struct {
	tx pgx.Tx
}

func NewTeamRepository(tx pgx.Tx) TeamRepository {
	return &postgresTeamRepository{tx: tx}
}

func (r *postgresTeamRepository) Create(ctx context.Context, team *Team) error {
	const query = `
		INSERT INTO teams (
			id,
			name,
			slug,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4,$5)
	`
	_, err := r.tx.Exec(ctx, query, team.ID, team.Name, team.Slug, team.CreatedAt, team.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create team: %w", err)
	}

	return nil
}

func (r *postgresTeamRepository) GetByID(ctx context.Context, id uuid.UUID) (*Team, error) {
	const query = `
		SELECT
			id,
			name,
			slug,
			created_at,
			updated_at
		FROM teams
		WHERE id = $1
	`
	var team Team
	err := r.tx.QueryRow(ctx, query, id).Scan(&team.ID, &team.Name, &team.Slug, &team.CreatedAt, &team.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get team by id: %w", err)
	}
	return &team, nil
}

func (r *postgresTeamRepository) GetBySlug(ctx context.Context, slug string) (*Team, error) {
	const query = `
		SELECT
			id,
			name,
			slug,
			created_at,
			updated_at
		FROM teams
		WHERE slug = $1
	`
	var team Team
	err := r.tx.QueryRow(ctx, query, slug).Scan(&team.ID, &team.Name, &team.Slug, &team.CreatedAt, &team.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get team by slug: %w", err)
	}
	return &team, nil
}

func (r *postgresTeamRepository) List(ctx context.Context) ([]*Team, error) {
	const query = `
		SELECT
			id,
			name,
			slug,
			created_at,
			updated_at
		FROM teams
		ORDER BY created_at,id
	`
	teams := make([]*Team, 0)

	rows, err := r.tx.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var team Team

		if err := rows.Scan(&team.ID, &team.Name, &team.Slug, &team.CreatedAt, &team.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan team: %w", err)
		}

		teams = append(teams, &team)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate teams: %w", err)
	}

	return teams, nil
}
