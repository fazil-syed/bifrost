package token

import (
	"errors"
	"fmt"
	"time"

	aero "github.com/aerospike/aerospike-client-go/v8"
	"github.com/aerospike/aerospike-client-go/v8/types"
	"github.com/google/uuid"
)

const tokenSet = "tokens"

const (
	binType        = "type"
	binTenantID    = "tenant_id"
	binApplication = "application"
	binAudience    = "audience"
	binUserID      = "user_id"
	binScopes      = "scopes"
	binIssuedAt    = "issued_at"
	binExpiresAt   = "expires_at"
	binRevokedAt   = "revoked_at"
	binFamilyID    = "family_id"
	binStatus      = "status"
	binUsedAt      = "used_at"
)

type AerospikeTokenRepository struct {
	client      *aero.Client
	namespace   string
	readPolicy  *aero.BasePolicy
	writePolicy *aero.WritePolicy
}

func NewAerospikeTokenRepository(
	client *aero.Client,
	namespace string,
	readPolicy *aero.BasePolicy,
	writePolicy *aero.WritePolicy,
) TokenRepository {
	return &AerospikeTokenRepository{
		client:      client,
		namespace:   namespace,
		readPolicy:  readPolicy,
		writePolicy: writePolicy,
	}

}

func (r *AerospikeTokenRepository) Create(token *Token) error {
	if token == nil {
		return fmt.Errorf("token is required")
	}

	key, err := r.key(token.ID)
	if err != nil {
		return err
	}

	policy := *r.writePolicy
	policy.RecordExistsAction = aero.CREATE_ONLY
	policy.Expiration = ttlSeconds(token.IssuedAt, token.ExpiresAt)

	bins := tokenBins(token)

	if err := r.client.PutBins(&policy, key, binsToBins(bins)...); err != nil {
		return fmt.Errorf("create token: %w", err)
	}
	return nil
}

func (r *AerospikeTokenRepository) GetByID(id string) (*Token, error) {
	if id == "" {
		return nil, fmt.Errorf("token ID is required")
	}

	key, err := r.key(id)
	if err != nil {
		return nil, err
	}

	record, err := r.client.Get(r.readPolicy, key)
	if err != nil {
		if errors.Is(err, aero.ErrKeyNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, fmt.Errorf("get token: %w", err)
	}

	token, err := tokenFromRecord(record)

	if err != nil {
		return nil, fmt.Errorf("decode token: %w", err)
	}
	return token, nil
}

func (r *AerospikeTokenRepository) Revoke(id string) error {
	if id == "" {
		return fmt.Errorf("token ID is required")
	}

	key, err := r.key(id)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	policy := *r.writePolicy
	policy.RecordExistsAction = aero.UPDATE_ONLY
	policy.Expiration = aero.TTLDontUpdate

	bins := []*aero.Bin{
		{
			Name:  binRevokedAt,
			Value: aero.NewLongValue(now.UnixNano()),
		},
	}
	if err := r.client.PutBins(&policy, key, bins...); err != nil {
		if errors.Is(err, aero.ErrKeyNotFound) {
			return ErrTokenNotFound
		}
		return fmt.Errorf("revoke token: %w", err)
	}

	return nil
}

func (r *AerospikeTokenRepository) ConsumeRefreshToken(id string, now time.Time) (*Token, error) {
	if id == "" {
		return nil, fmt.Errorf("token ID is required")
	}

	key, err := r.key(id)
	if err != nil {
		return nil, err
	}

	record, err := r.client.Get(r.readPolicy, key)

	if err != nil {
		if errors.Is(err, aero.ErrKeyNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, fmt.Errorf("get refresh token: %w", err)
	}
	token, err := tokenFromRecord(record)
	if err != nil {
		return nil, fmt.Errorf("decode refresh token: %w", err)
	}

	if token.Type != TypeRefresh {
		return nil, fmt.Errorf("token %q is not a refresh token", id)
	}

	if !token.ExpiresAt.After(now) {
		return nil, ErrTokenExpired
	}

	if token.RevokedAt != nil {
		return nil, ErrTokenRevoked
	}
	if token.Status == RefreshTokenUsed {
		return nil, ErrRefreshTokenUsed
	}
	if token.Status == RefreshTokenRevoked {
		return nil, ErrTokenRevoked
	}

	// Compare and set using record generation
	// Only the request holding this generation may transition ACTIVE-USED

	policy := *r.writePolicy
	policy.RecordExistsAction = aero.UPDATE_ONLY
	policy.GenerationPolicy = aero.EXPECT_GEN_EQUAL
	policy.Generation = record.Generation
	policy.Expiration = aero.TTLDontUpdate

	bins := []*aero.Bin{
		{
			Name:  binStatus,
			Value: aero.NewStringValue(string(RefreshTokenUsed)),
		},
		{
			Name:  binUsedAt,
			Value: aero.NewLongValue(now.UnixNano()),
		},
	}

	if err := r.client.PutBins(&policy, key, bins...); err != nil {
		if err.Matches(types.GENERATION_ERROR) {
			// Another request changed this token after our read
			// Re-read so we can return a meaningful domain error
			current, readErr := r.GetByID(id)
			if readErr != nil {
				return nil, readErr
			}
			if current.Status == RefreshTokenUsed {
				return nil, ErrRefreshTokenUsed
			}
			if current.RevokedAt != nil || current.Status == RefreshTokenRevoked {
				return nil, ErrTokenRevoked
			}
			return nil, fmt.Errorf("refresh token changed concurrently")
		}
		return nil, fmt.Errorf("consume refresh token: %w", err)
	}
	token.UsedAt = timePtr(now)
	token.Status = RefreshTokenUsed
	return token, nil
}

func (r *AerospikeTokenRepository) key(id string) (*aero.Key, error) {
	if id == "" {
		return nil, fmt.Errorf("token ID is required")
	}

	key, err := aero.NewKey(r.namespace, tokenSet, id)
	if err != nil {
		return nil, fmt.Errorf("create token key: %w", err)
	}
	return key, nil
}

func ttlSeconds(issuedAt, expiresAt time.Time) uint32 {
	duration := expiresAt.Sub(issuedAt)
	if duration <= 0 {
		return 1
	}
	seconds := uint64(duration / time.Second)
	if duration%time.Second != 0 {
		seconds++
	}
	if seconds == 0 {
		seconds = 1
	}

	return uint32(seconds)
}

func timePtr(t time.Time) *time.Time {
	value := t
	return &value
}

func binsToBins(binMap aero.BinMap) []*aero.Bin {
	bins := make([]*aero.Bin, 0, len(binMap))

	for name, value := range binMap {
		bins = append(bins, aero.NewBin(name, value))
	}
	return bins
}

func tokenBins(token *Token) aero.BinMap {
	bins := aero.BinMap{
		binType:        string(token.Type),
		binTenantID:    token.TenantID.String(),
		binApplication: token.ApplicationID.String(),
		binAudience:    token.Audience,
		binUserID:      token.UserID.String(),
		binIssuedAt:    token.IssuedAt.UnixNano(),
		binExpiresAt:   token.ExpiresAt.UnixNano(),
	}

	if len(token.Scopes) > 0 {
		bins[binScopes] = token.Scopes
	}
	if token.RevokedAt != nil {
		bins[binRevokedAt] = token.RevokedAt.UnixNano()
	}

	if token.Type == TypeRefresh {
		bins[binFamilyID] = token.FamilyID.String()
		bins[binStatus] = string(token.Status)

		if token.UsedAt != nil {
			bins[binUsedAt] = token.UsedAt.UnixNano()
		}
	}
	return bins
}

func tokenFromRecord(record *aero.Record) (*Token, error) {
	typeValue, ok := record.Bins[binType].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token type")
	}
	tenantID, err := parseUUIDBin(record.Bins[binTenantID], binTenantID)

	if err != nil {
		return nil, err
	}
	applicationID, err := parseUUIDBin(record.Bins[binApplication], binApplication)

	if err != nil {
		return nil, err
	}
	userID, err := parseUUIDBin(record.Bins[binUserID], binUserID)

	if err != nil {
		return nil, err
	}

	audience, ok := record.Bins[binAudience].(string)

	if !ok || audience == "" {
		return nil, fmt.Errorf("invalid audience")
	}

	issuedAt, err := parseTimeBin(record.Bins[binIssuedAt], binIssuedAt)

	if err != nil {
		return nil, err
	}

	expiresAt, err := parseTimeBin(record.Bins[binExpiresAt], binExpiresAt)

	if err != nil {
		return nil, err
	}

	token := &Token{
		ID:            record.Key.Value().String(),
		Type:          Type(typeValue),
		TenantID:      tenantID,
		ApplicationID: applicationID,
		UserID:        userID,
		Audience:      audience,
		IssuedAt:      issuedAt,
		ExpiresAt:     expiresAt,
		Scopes:        parseScopes(record.Bins[binScopes]),
	}

	if _, exists := record.Bins[binRevokedAt]; exists {
		revokedAt, err := parseTimeBin(record.Bins[binRevokedAt], binRevokedAt)
		if err != nil {
			return nil, err
		}
		token.RevokedAt = &revokedAt
	}

	if token.Type == TypeRefresh {
		familyID, err := parseUUIDBin(record.Bins[binFamilyID], binFamilyID)
		if err != nil {
			return nil, err
		}
		status, ok := record.Bins[binStatus].(string)

		if !ok {
			return nil, fmt.Errorf("invalid refresh token status")
		}

		token.FamilyID = familyID
		token.Status = RefreshTokenStatus(status)

		if value, exists := record.Bins[binUsedAt]; exists {
			usedAt, err := parseTimeBin(value, binUsedAt)

			if err != nil {
				return nil, err
			}
			token.UsedAt = &usedAt
		}
	}
	return token, nil
}

func parseScopes(value any) []string {
	switch scopes := value.(type) {
	case []string:
		return cloneScoppes(scopes)
	case []interface{}:
		result := make([]string, 0, len(scopes))
		for _, scope := range scopes {
			if value, ok := scope.(string); ok {
				result = append(result, value)
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	default:
		return nil
	}

}

func parseTimeBin(value any, name string) (time.Time, error) {
	raw, ok := value.(int64)
	if !ok {
		return time.Time{}, fmt.Errorf("invalid %s", name)
	}
	return time.Unix(0, raw).UTC(), nil
}

func parseUUIDBin(value any, name string) (uuid.UUID, error) {
	raw, ok := value.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid %s", name)
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse %s: %w", name, err)
	}
	return id, nil
}
