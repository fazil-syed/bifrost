package authorization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type service struct {
	db *pgxpool.Pool
}

func (s *service) GetScopes(ctx context.Context, userID uuid.UUID, applicationID uuid.UUID) ([]string, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})

	if err != nil {
		return nil, fmt.Errorf("begin get scopes transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewAuthorizationRepository(tx)

	return repository.GetScopes(ctx, userID, applicationID)

}
