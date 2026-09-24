package role

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID            uuid.UUID
	ApplicationID uuid.UUID
	Name          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func New(applicationID uuid.UUID, name string, now time.Time) (*Role, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, ErrRoleNameInvalid
	}

	id, err := uuid.NewV7()

	if err != nil {
		return nil, fmt.Errorf("generate role id: %w", err)
	}
	return &Role{
		ID:            id,
		ApplicationID: applicationID,
		Name:          name,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}
