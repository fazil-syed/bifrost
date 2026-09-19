package tenant

import "errors"

var (
	ErrTenantNotFound = errors.New("tenant not found")
	ErrTenantDisabled = errors.New("tenant disabled")

	ErrTenantSlugExists         = errors.New("tenant slug already exists")
	ErrTenantDatabaseNameExists = errors.New("tenant database name already exists")
)
