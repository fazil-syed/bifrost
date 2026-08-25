package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *User) error {
	const query = `
		INSERT INTO users (
			id,
			status,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4)
	`

	_, err := r.db.Exec(ctx, query, user.ID, user.Status, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	const query = `
		SELECT 
			id,
			status,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var user User

	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user %s: %w", id, err)
	}
	return &user, nil

}

func (r *postgresUserRepository) Disable(ctx context.Context, id uuid.UUID, updatedAt time.Time) error {
	const query = `
		UPDATE users
		SET 
			status = $1,
			updated_at = $2
		WHERE id = $3
	`
	result, err := r.db.Exec(
		ctx,
		query,
		UserStatusDisabled,
		updatedAt,
		id,
	)
	if err != nil {
		return fmt.Errorf("disable user %s: %w", id, err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *postgresUserRepository) Enable(ctx context.Context, id uuid.UUID, updatedAt time.Time) error {
	const query = `
		UPDATE users
		SET 
			status = $1,
			updated_at = $2
		WHERE id = $3
	`
	result, err := r.db.Exec(
		ctx,
		query,
		UserStatusActive,
		updatedAt,
		id,
	)
	if err != nil {
		return fmt.Errorf("enable user %s: %w", id, err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}
