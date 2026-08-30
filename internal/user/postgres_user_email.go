package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postgresUserEmailRepository struct {
	tx pgx.Tx
}

func NewPostgresUserEmailRepository(tx pgx.Tx) UserEmailRepository {
	return &postgresUserEmailRepository{
		tx: tx,
	}
}

func (r *postgresUserEmailRepository) Create(ctx context.Context, userEmail *UserEmail) error {
	const query = `
		INSERT INTO user_emails (
			user_id,
			email,
			verified,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5)
	`

	_, err := r.tx.Exec(ctx, query, userEmail.UserID, userEmail.Email, userEmail.Verified, userEmail.CreatedAt, userEmail.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create user email: %w", err)
	}
	return nil
}

func (r *postgresUserEmailRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*UserEmail, error) {
	const query = `
		SELECT
			user_id,
			email,
			verified,
			created_at,
			updated_at
		FROM user_emails
		WHERE user_id = $1
	`
	var userEmail UserEmail
	err := r.tx.QueryRow(ctx, query, userID).Scan(&userEmail.UserID, &userEmail.Email, &userEmail.Verified, &userEmail.CreatedAt, &userEmail.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf(
			"get user email for user %s: %w",
			userID,
			err,
		)
	}
	return &userEmail, nil
}

func (r *postgresUserEmailRepository) GetByEmail(
	ctx context.Context,
	email Email,
) (*UserEmail, error) {
	const query = `
		SELECT
			user_id,
			email,
			verified,
			created_at,
			updated_at
		FROM user_emails
		WHERE email = $1
	`
	var userEmail UserEmail
	err := r.tx.QueryRow(ctx, query, email).Scan(&userEmail.UserID, &userEmail.Email, &userEmail.Verified, &userEmail.CreatedAt, &userEmail.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserEmailNotFound
		}
		return nil, fmt.Errorf(
			"get user email for user %s: %w",
			email,
			err,
		)
	}
	return &userEmail, nil
}
