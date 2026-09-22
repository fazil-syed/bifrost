package team

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MembershipType string

const (
	MembershipMember MembershipType = "MEMBER"
	MembershipOwner  MembershipType = "OWNER"
)

type Membership struct {
	TeamID         uuid.UUID
	UserID         uuid.UUID
	MembershipType MembershipType
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewMembership(teamID uuid.UUID, userID uuid.UUID, membershipType MembershipType, now time.Time) (*Membership, error) {
	if teamID == uuid.Nil {
		return nil, fmt.Errorf("team ID is required")
	}

	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	switch membershipType {
	case MembershipMember, MembershipOwner:
	default:
		return nil, fmt.Errorf("invalid membership type %q", membershipType)
	}

	return &Membership{
		TeamID:         teamID,
		UserID:         userID,
		MembershipType: membershipType,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}
