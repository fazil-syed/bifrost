package authorization

import (
	"context"

	"github.com/google/uuid"
)

type AuthorizationService interface {
	GetScopes(
		ctx context.Context,
		userID uuid.UUID,
		applicationID uuid.UUID,
	) ([]string, error)
}
