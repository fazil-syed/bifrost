package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postgresUserRepository struct {
	tx pgx.Tx
}

func NewPostgresUserRepository(tx pgx.Tx) UserRepository {
	return &postgresUserRepository{tx: tx}
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

	_, err := r.tx.Exec(ctx, query, user.ID, user.Status, user.CreatedAt, user.UpdatedAt)
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

	err := r.tx.QueryRow(ctx, query, id).Scan(&user.ID, &user.Status, &user.CreatedAt, &user.UpdatedAt)

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
	result, err := r.tx.Exec(
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
	result, err := r.tx.Exec(
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
func (r *postgresUserRepository) GetByEmail(ctx context.Context, email Email) (*User, error) {
	const query = `
		SELECT
			u.id,
			u.status,
			u.created_at,
			u.updated_at
		FROM users u
		INNER JOIN user_emails ue
			ON ue.user_id = u.id
		WHERE ue.email = $1
	`

	var user User
	err := r.tx.QueryRow(ctx, query, email).Scan(&user.ID, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email %s: %w", email, err)
	}
	return &user, nil
}
