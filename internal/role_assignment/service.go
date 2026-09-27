package roleassignment

import (
	"context"

	"github.com/google/uuid"
)

type RoleAssignmentService interface {
	AssignToUser(ctx context.Context, roleID uuid.UUID, userID uuid.UUID) (*RoleAssignment, error)

	AssignToTeam(ctx context.Context, roleID uuid.UUID, teamID uuid.UUID) (*RoleAssignment, error)

	GetByID(ctx context.Context, id uuid.UUID) (*RoleAssignment, error)

	ListByRole(ctx context.Context, roleID uuid.UUID) ([]*RoleAssignment, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*RoleAssignment, error)
	ListByTeam(ctx context.Context, teamID uuid.UUID) ([]*RoleAssignment, error)

	Remove(ctx context.Context, id uuid.UUID) error
}
