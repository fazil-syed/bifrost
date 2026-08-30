package user

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ExternalIdentity struct {
	UserID    uuid.UUID
	Issuer    string
	Subject   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrExternalIdentityNotFound = errors.New("external identity not found ")

func NewExternalIdentity(
	userID uuid.UUID,
	issuer string,
	subject string,
	now time.Time,
) (*ExternalIdentity, error) {

	issuer, err := normalizeIssuer(issuer)
	if err != nil {
		return nil, err
	}
	if subject == "" {
		return nil, fmt.Errorf("external identity subject is required")
	}
	return &ExternalIdentity{
		UserID:    userID,
		Issuer:    issuer,
		Subject:   subject,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func normalizeIssuer(value string) (string, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return "", fmt.Errorf("external identity issuer is required")
	}

	parsed, err := url.Parse(value)

	if err != nil {
		return "", fmt.Errorf("invalid external identity issuer: %w", err)
	}

	if parsed.Scheme != "https" {
		return "", fmt.Errorf("external identity issuer must use https")
	}

	if parsed.Host == "" {
		return "", fmt.Errorf("external identity issuer must contain a host")
	}
	if parsed.RawQuery != "" {
		return "", fmt.Errorf("external identity issuer must not contain a query")
	}

	if parsed.Fragment != "" {
		return "", fmt.Errorf("external identity issuer must not contain a fragment")
	}

	return value, nil
}
