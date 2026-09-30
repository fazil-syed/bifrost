package authorization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postgresAuthorizationRepository struct {
	tx pgx.Tx
}

func NewAuthorizationRepository(tx pgx.Tx) AuthorizationRepository {
	return &postgresAuthorizationRepository{tx: tx}
}

func (r *postgresAuthorizationRepository) GetScopes(ctx context.Context, userID uuid.UUID, applicationID uuid.UUID) ([]string, error) {
	const query = `
		SELECT DISTINCT r.name
		FROM roles r
		JOIN role_assignments ra
			ON ra.role_id = r.id
		LEFT JOIN team_memberships tm
			ON tm.team_id = ra.team_id
			AND tm.user_id = $1
		WHERE r.application_id = $2
			AND (
				ra.user_id = $1
				OR tm.user_id IS NOT NULL
			)
		ORDER BY r.name
	`

	rows, err := r.tx.Query(ctx, query, userID, applicationID)

	if err != nil {
		return nil, fmt.Errorf("get roles assigned to user %s for application %s: %w", userID, applicationID, err)
	}

	defer rows.Close()

	scopes := make([]string, 0)
	for rows.Next() {
		var scope string

		if err := rows.Scan(&scope); err != nil {
			return nil, fmt.Errorf("scan scope: %w", err)
		}
		scopes = append(scopes, scope)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate scopes: %w", err)
	}
	return scopes, err
}
