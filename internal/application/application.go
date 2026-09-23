package application

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	OwnerUserID *uuid.UUID
	OwnerTeamID *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func New(name string, slug string, ownerUserID *uuid.UUID, ownerTeamID *uuid.UUID, now time.Time) (*Application, error) {
	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(slug)

	if name == "" {
		return nil, ErrApplicationNameInvalid
	}
	if slug == "" {
		return nil, ErrApplicationSlugInvalid
	}

	if (ownerTeamID == nil) == (ownerUserID == nil) {
		return nil, ErrApplicationOwnerInvalid
	}
	id, err := uuid.NewV7()

	if err != nil {
		return nil, fmt.Errorf("generate application id: %w", err)
	}

	return &Application{
		ID: id,
		Name: name,
		Slug: slug,
		OwnerUserID: ownerUserID,
		OwnerTeamID: ownerTeamID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
