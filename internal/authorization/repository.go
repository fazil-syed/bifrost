package authorization

import (
	"context"

	"github.com/google/uuid"
)

type AuthorizationRepository interface {
	GetScopes(
		ctx context.Context,
		userID uuid.UUID,
		applicationID uuid.UUID,
	) ([]string, error)
}
