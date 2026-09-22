package team

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MembershipRepository interface {
	Create(ctx context.Context, membership *Membership) error

	Get(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (*Membership, error)

	ListByTeam(ctx context.Context, teamID uuid.UUID) ([]*Membership, error)

	UpdateType(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, membershipType MembershipType, updatedAt time.Time) error

	Delete(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error

	CountOwners(ctx context.Context, teamID uuid.UUID) (int, error)
}
