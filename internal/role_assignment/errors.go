package roleassignment

import "errors"

var (
	ErrRoleAssignmentNotFound = errors.New("role assignment not found")

	ErrRoleAssingmentSubjectInvalid = errors.New("role assignment must have exactly one subject")

	ErrRoleAssignmentAlreadyExists = errors.New("role assignment already exists")

	ErrRoleAssignmentUnauthorized = errors.New("user is not authorized to manage this role assignment")
)
