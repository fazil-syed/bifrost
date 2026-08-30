package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postgresPasswordCredentialRepository struct {
	tx pgx.Tx
}

func NewPostgresPasswordCredentialRepository(tx pgx.Tx) PasswordCredentialRepository {
	return &postgresPasswordCredentialRepository{tx: tx}
}

func (r *postgresPasswordCredentialRepository) Create(ctx context.Context, credential *PasswordCredential) error {
	const query = `
		INSERT INTO credentials (
			user_id,
			password_hash,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4)
	`

	_, err := r.tx.Exec(ctx, query, credential.UserID, credential.PasswordHash, credential.CreatedAt, credential.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create password credential: %w", err)
	}
	return nil

}

func (r *postgresPasswordCredentialRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*PasswordCredential, error) {
	const query = `
		SELECT 
			user_id,
			password_hash,
			created_at,
			updated_at
		FROM credentials
		WHERE user_id = $1
	`

	var passwordCredential PasswordCredential

	err := r.tx.QueryRow(ctx, query, userID).Scan(&passwordCredential.UserID, passwordCredential.PasswordHash, passwordCredential.CreatedAt, passwordCredential.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPasswordCredentialNotFound
		}
		return nil, fmt.Errorf("get password credential %s : %w", userID, err)
	}
	return &passwordCredential, nil
}

func (r *postgresPasswordCredentialRepository) Update(ctx context.Context, credential *PasswordCredential) error {
	const query = `
		UPDATE credentials
		SET 
			password_hash = $1,
			updated_at = $2
		WHERE user_id = $3
	`

	result, err := r.tx.Exec(ctx, query, credential.PasswordHash, credential.UpdatedAt, credential.UserID)

	if err != nil {
		return fmt.Errorf("updated password credential for user %s: %w", credential.UserID, err)
	}

	if result.RowsAffected() == 0 {
		return ErrPasswordCredentialNotFound
	}
	return nil

}
