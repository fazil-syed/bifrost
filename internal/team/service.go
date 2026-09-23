package team

import (
	"context"

	"github.com/google/uuid"
)

type TeamService interface {
	Create(ctx context.Context, name string, slug string, ownerUserID uuid.UUID) (*Team, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Team, error)
	GetBySlug(ctx context.Context, slug string) (*Team, error)
	List(ctx context.Context) ([]*Team, error)
}

type MembershipService interface {
	AddMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error

	GetMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (*Membership, error)

	ListMembers(ctx context.Context, teamID uuid.UUID) ([]*Membership, error)

	PromoteToOwner(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error

	DemoteOwner(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error

	RemoveMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error
}
