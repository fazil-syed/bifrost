package user

import (
	"context"

	"github.com/google/uuid"
)

type ExternalIdentityRepository interface {
	Create(
		ctx context.Context,
		externalIdentity ExternalIdentity,
	) error

	GetByIssuerAndSubject(
		ctx context.Context,
		issuer string,
		subject string,
	) (*ExternalIdentity, error)

	GetByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]*ExternalIdentity, error)
}
