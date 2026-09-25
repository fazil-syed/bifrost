package roleassignment

import (
	"context"

	"github.com/google/uuid"
)

type RoleAssignmentRepository interface {
	Create(ctx context.Context, assignment *RoleAssignment) error

	GetByID(ctx context.Context, id uuid.UUID) (*RoleAssignment, error)

	GetByRoleAndUser(ctx context.Context, roleID uuid.UUID, userID uuid.UUID) (*RoleAssignment, error)

	GetByRoleAndTeam(ctx context.Context, roleID uuid.UUID, teamID uuid.UUID) (*RoleAssignment, error)

	ListByRole(ctx context.Context, roleID uuid.UUID) ([]*RoleAssignment, error)

	ListByUser(ctx context.Context, userID uuid.UUID) ([]*RoleAssignment, error)

	ListByTeam(ctx context.Context, teamID uuid.UUID) ([]*RoleAssignment, error)

	Delete(ctx context.Context, id uuid.UUID) error
}
