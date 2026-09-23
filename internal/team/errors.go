package team

import "errors"

var (
	ErrTeamNotFound           = errors.New("team not found")
	ErrTeamSlugExists         = errors.New("team slug already exists")
	ErrTeamMembershipNotFound = errors.New("team membership not found")
	ErrTeamMembershipExists   = errors.New("team membership already exists")
	ErrInvalidMembershipType  = errors.New("invalid membership type")
	ErrTeamMustHaveOwner      = errors.New("team must have at least one owner")
)
