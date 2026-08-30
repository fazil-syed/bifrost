package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postgresExternalIdentityRepository struct {
	tx pgx.Tx
}

func NewPostgresExternalIdentityRepository(tx pgx.Tx) ExternalIdentityRepository {
	return &postgresExternalIdentityRepository{
		tx: tx,
	}
}

func (r *postgresExternalIdentityRepository) Create(ctx context.Context, externalIdentity ExternalIdentity) error {
	const query = `
		INSERT INTO external_identities (
			user_id,
			issuer,
			subject,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5)
	`
	_, err := r.tx.Exec(ctx, query, externalIdentity.UserID, externalIdentity.Issuer, externalIdentity.Subject, externalIdentity.CreatedAt, externalIdentity.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create external identity : %w", err)
	}
	return nil

}

func (r *postgresExternalIdentityRepository) GetByIssuerAndSubject(ctx context.Context, issuer string, subject string) (*ExternalIdentity, error) {
	const query = `
		SELECT
			user_id,
			issuer,
			subject,
			created_at,
			updated_at
		FROM external_identities
		WHERE issuer = $1
			AND subject = $2
	`

	var externalIdentity ExternalIdentity

	err := r.tx.QueryRow(ctx, query, issuer, subject).Scan(&externalIdentity.UserID, &externalIdentity.Issuer, &externalIdentity.Subject, &externalIdentity.CreatedAt, &externalIdentity.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrExternalIdentityNotFound
		}
		return nil, fmt.Errorf("get external identity %s/%s: %w", issuer, subject, err)
	}
	return &externalIdentity, nil
}

func (r *postgresExternalIdentityRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*ExternalIdentity, error) {
	const query = `
		SELECT
			user_id,
			issuer,
			subject,
			created_at,
			updated_at
		FROM external_identities
		WHERE user_id = $1
		ORDER BY created_at
	`

	rows, err := r.tx.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get external identities for user %s: %w", userID, err)
	}

	defer rows.Close()

	var identities []*ExternalIdentity

	for rows.Next() {
		var identity ExternalIdentity

		if err := rows.Scan(
			&identity.UserID,
			&identity.Issuer,
			&identity.Subject,
			&identity.CreatedAt,
			&identity.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan external identity for user %s: %w", userID, err)
		}
		identities = append(identities, &identity)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate external identities for user %s: %w", userID, err)
	}

	return identities, nil
}
