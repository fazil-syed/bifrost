package user

import (
	"fmt"
	"net/mail"
	"strings"
)

type Email string

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)

	if value == "" {
		return "", fmt.Errorf("email is required")
	}
	address, err := mail.ParseAddress(value)
	if err != nil {
		return "", fmt.Errorf("invalid email: %w", err)
	}
	if address.Address != value {
		return "", fmt.Errorf("invalid email format")
	}

	return Email(value), nil
}

func (e Email) String() string {
	return string(e)
}
