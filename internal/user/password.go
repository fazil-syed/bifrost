package user

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

var passwordHashParams = &argon2id.Params{
	Memory:      20 * 1024,
	Iterations:  4,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password is required")
	}
	hash, err := argon2id.CreateHash(password, passwordHashParams)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return hash, nil
}

func VerifyPassword(password, hash string) (bool, error) {
	if password == "" {
		return false, fmt.Errorf("password is required")
	}

	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("verify password: %w", err)
	}
	return match, nil
}
