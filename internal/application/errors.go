package application

import "errors"

var (
	ErrApplicationNotFound   = errors.New("application not found")
	ErrApplicationSlugExists = errors.New("application slug already exists")

	ErrApplicationNameInvalid = errors.New("application name is invalid")

	ErrApplicationSlugInvalid = errors.New("application slug is invalid")

	ErrApplicationOwnerInvalid = errors.New("application must have exactly one owner")
)
