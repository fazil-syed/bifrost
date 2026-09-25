package roleassignment

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RoleAssignment struct {
	ID        uuid.UUID
	RoleID    uuid.UUID
	UserID    *uuid.UUID
	TeamID    *uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(roleID uuid.UUID, userID *uuid.UUID, teamID *uuid.UUID, now time.Time) (*RoleAssignment, error) {
	if (userID == nil) == (teamID == nil) {
		return nil, ErrRoleAssingmentSubjectInvalid
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("create role assignment id: %w", err)
	}

	return &RoleAssignment{
		ID:        id,
		RoleID:    roleID,
		UserID:    userID,
		TeamID:    teamID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
