package team

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(name, slug string, now time.Time) (*Team, error) {
	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(slug)
	if name == "" {
		return nil, fmt.Errorf("team name is required")
	}

	if slug == "" {
		return nil, fmt.Errorf("team slug is required")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate team ID: %w", err)
	}

	return &Team{
		ID:        id,
		Name:      name,
		Slug:      slug,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
