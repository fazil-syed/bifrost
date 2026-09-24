package role

import "errors"

var (
	ErrRoleNotFound = errors.New("role not found")

	ErrRoleNameInvalid = errors.New("role name is invalid")

	ErrRoleExists = errors.New("role already exists")

	ErrRoleUnauthorized = errors.New("user is not authorized to manager this role")
)
