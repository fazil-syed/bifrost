package tenant

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusDisabled Status = "DISABLED"
)

type Tenant struct {
	ID           uuid.UUID
	Name         string
	Slug         string
	DatabaseName string
	Status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func New(name, slug, databaseName string, now time.Time) (*Tenant, error) {
	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(slug)
	databaseName = strings.TrimSpace(databaseName)
	if name == "" {
		return nil, fmt.Errorf("tenant name is required")
	}

	if slug == "" {
		return nil, fmt.Errorf("tenant slug is required")
	}

	if databaseName == "" {
		return nil, fmt.Errorf("tenant database name is required")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate tenant ID: %w", err)
	}
	return &Tenant{
		ID:           id,
		Name:         name,
		Slug:         slug,
		DatabaseName: databaseName,
		Status:       StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
