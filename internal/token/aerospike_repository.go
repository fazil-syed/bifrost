package token

import (
	"errors"
	"fmt"
	"time"

	aero "github.com/aerospike/aerospike-client-go/v8"
	"github.com/google/uuid"
)

const tokenSet = "tokens"
const tokenFamilySet = "token_families"

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

const (
	binFamilyCreatedAt = "created_at"
	binFamilyRevokedAt = "revoked_at"
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

func (r *AerospikeTokenRepository) CreateInitialTokenPair(family *RefreshTokenFamily, accessToken *Token, refreshToken *Token) error {
	if family == nil {
		return fmt.Errorf("refresh token family is required")
	}
	if accessToken == nil {
		return fmt.Errorf("access token is required")
	}
	if refreshToken == nil {
		return fmt.Errorf("refresh token is required")
	}

	if refreshToken.Type != TypeRefresh {
		return fmt.Errorf("refresh token has invalid type")
	}

	if refreshToken.FamilyID != family.ID {
		return fmt.Errorf("refresh token does not belong to family")
	}

	txn := aero.NewTxn()

	familyKey, err := aero.NewKey(r.namespace, tokenFamilySet, family.ID.String())

	if err != nil {
		return fmt.Errorf("create refresh token family key: %w", err)
	}

	accessKey, err := aero.NewKey(r.namespace, tokenSet, accessToken.ID)

	if err != nil {
		return fmt.Errorf("create access token key: %w", err)
	}
	refreshKey, err := aero.NewKey(r.namespace, tokenSet, refreshToken.ID)
	if err != nil {
		return fmt.Errorf("create refresh token key: %w", err)
	}

	writePolicy := *r.writePolicy

	writePolicy.Txn = txn
	writePolicy.RecordExistsAction = aero.CREATE_ONLY

	if err := r.client.PutBins(&writePolicy, familyKey, aero.NewBin(binFamilyCreatedAt, family.CreatedAt.UnixNano())); err != nil {
		r.client.Abort(txn)
		return fmt.Errorf("create token family: %w", err)
	}
	if err := r.client.PutBins(&writePolicy, accessKey, binsToBins(tokenBins(accessToken))...); err != nil {
		r.client.Abort(txn)
		return fmt.Errorf("create access token: %w", err)
	}
	if err := r.client.PutBins(&writePolicy, refreshKey, binsToBins(tokenBins(refreshToken))...); err != nil {
		r.client.Abort(txn)
		return fmt.Errorf("create refresh token: %w", err)
	}

	if _, err := r.client.Commit(txn); err != nil {
		return fmt.Errorf("commit initial token pair: %w", err)
	}

	return nil

}

func (r *AerospikeTokenRepository) RotateRefreshToken(familyID string, presentedRefreshTokenID string, accessToken *Token, refreshToken *Token, now time.Time) error {
	if familyID == "" {
		return fmt.Errorf("refresh token family id is required")
	}

	if presentedRefreshTokenID == "" {
		return fmt.Errorf("presented refresh token ID is required")
	}

	if accessToken == nil {
		return fmt.Errorf("access token is required")
	}
	if refreshToken == nil {
		return fmt.Errorf("refresh token is required")
	}

	if refreshToken.Type != TypeRefresh {
		return fmt.Errorf("refresh token has invalid type")
	}

	fID, err := uuid.Parse(familyID)
	if err != nil {
		return fmt.Errorf("parse refresh token family id: %w", err)
	}
	if refreshToken.FamilyID != fID {
		return fmt.Errorf("replacement refresh token does not belong to family")
	}
	familyKey, err := aero.NewKey(r.namespace, tokenFamilySet, familyID)

	if err != nil {
		return fmt.Errorf("create refresh token family key: %w", err)
	}

	presentedKey, err := aero.NewKey(r.namespace, tokenSet, presentedRefreshTokenID)

	if err != nil {
		return fmt.Errorf("create presented refresh token key: %w", err)
	}

	replacmentAccessKey, err := aero.NewKey(r.namespace, tokenSet, accessToken.ID)

	if err != nil {
		return fmt.Errorf("create replacement access token key: %w", err)
	}

	replacementRefreshKey, err := aero.NewKey(r.namespace, tokenSet, refreshToken.ID)

	if err != nil {
		return fmt.Errorf("create replacement refresh token key: %w", err)
	}

	txn := aero.NewTxn()

	readPolicy := *r.readPolicy
	readPolicy.Txn = txn

	writePolicy := *r.writePolicy

	writePolicy.Txn = txn

	familyRecord, err := r.client.Get(&readPolicy, familyKey)

	if err != nil {
		r.client.Abort(txn)
		if errors.Is(err, aero.ErrKeyNotFound) {
			return ErrTokenNotFound
		}
		return fmt.Errorf("read refresh token family: %w", err)
	}

	presentedRecord, err := r.client.Get(&readPolicy, presentedKey)

	if err != nil {
		r.client.Abort(txn)
		if errors.Is(err, aero.ErrKeyNotFound) {
			return ErrTokenNotFound
		}
		return fmt.Errorf("read presented refresh token: %w", err)
	}

	familyRevoked := false

	if value, exists := familyRecord.Bins[binFamilyRevokedAt]; exists {
		familyRevoked = value != nil
	}

	if familyRevoked {
		r.client.Abort(txn)
		return ErrTokenRevoked
	}

	presented, err := tokenFromRecord(presentedRecord)

	if err != nil {
		r.client.Abort(txn)
		return fmt.Errorf("decode presented refresh token: %w", err)
	}

	if presented.Type != TypeRefresh {
		r.client.Abort(txn)
		return ErrInvalidTokenType
	}

	if presented.FamilyID.String() != familyID {
		r.client.Abort(txn)
		return ErrTokenNotFound
	}

	if presented.RevokedAt != nil {
		r.client.Abort(txn)
		return ErrTokenRevoked
	}

	if !presented.ExpiresAt.After(now) {
		r.client.Abort(txn)
		return ErrTokenExpired
	}

	// Reuse of already consumed RT is the replay attach path

	if presented.Status == RefreshTokenUsed {
		writePolicy.RecordExistsAction = aero.UPDATE_ONLY
		if err := r.client.PutBins(&writePolicy, familyKey, aero.NewBin(binFamilyRevokedAt, now.UnixNano())); err != nil {
			r.client.Abort(txn)
			return fmt.Errorf("revoke refresh token family after replay: %w", err)
		}

		if _, err := r.client.Commit(txn); err != nil {
			return fmt.Errorf("commit refresh token replay revocation: %w", err)
		}
		return ErrRefreshTokenReplay
	}

	if presented.Status != RefreshTokenActive {
		r.client.Abort(txn)
		return ErrTokenRevoked
	}

	// normal rotation

	writePolicy.RecordExistsAction = aero.UPDATE_ONLY
	//consume current refresh token
	if err := r.client.PutBins(&writePolicy, presentedKey, aero.NewBin(binStatus, string(RefreshTokenUsed)), aero.NewBin(binUsedAt, now.UnixNano())); err != nil {
		r.client.Abort(txn)
		return fmt.Errorf("consume refresh token: %w", err)
	}

	writePolicy.RecordExistsAction = aero.CREATE_ONLY

	//create access token
	if err := r.client.PutBins(&writePolicy, replacmentAccessKey, binsToBins(tokenBins(accessToken))...); err != nil {
		r.client.Abort(txn)
		return fmt.Errorf("create replacement access token: %w", err)
	}

	//create refresh token

	if err := r.client.PutBins(&writePolicy, replacementRefreshKey, binsToBins(tokenBins(refreshToken))...); err != nil {
		return fmt.Errorf("create replacement refresh token: %w", err)
	}

	if _, err := r.client.Commit(txn); err != nil {
		return fmt.Errorf("commit refresh token rotation: %w", err)
	}
	return nil
}

func (r *AerospikeTokenRepository) CreateToken(token *Token) error {
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

func (r *AerospikeTokenRepository) GetTokenByID(id string) (*Token, error) {
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

func (r *AerospikeTokenRepository) RevokeToken(id string) error {
	if id == "" {
		return fmt.Errorf("token ID is required")
	}

	key, err := r.key(id)
	if err != nil {
		return err
	}

	now := time.Now()

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

func (r *AerospikeTokenRepository) CreateRefreshTokenFamily(family *RefreshTokenFamily) error {
	if family == nil {
		return fmt.Errorf("refresh token family is required")
	}

	key, err := aero.NewKey(r.namespace, tokenFamilySet, family.ID.String())
	if err != nil {
		return fmt.Errorf("create refresh token family key: %w", err)
	}

	policy := *r.writePolicy
	policy.RecordExistsAction = aero.CREATE_ONLY
	policy.Expiration = 0

	bins := aero.BinMap{
		binFamilyCreatedAt: family.CreatedAt.UnixNano(),
	}

	if family.RevokedAt != nil {
		bins[binFamilyRevokedAt] = family.RevokedAt.UnixNano()
	}

	if err := r.client.PutBins(&policy, key, binsToBins(bins)...); err != nil {
		return fmt.Errorf("create refresh token family: %w", err)
	}

	return nil

}

func (r *AerospikeTokenRepository) GetRefreshTokenFamilyById(id string) (*RefreshTokenFamily, error) {
	if id == "" {
		return nil, fmt.Errorf("refresh token family id is required")
	}

	key, err := aero.NewKey(r.namespace, tokenFamilySet, id)
	if err != nil {
		return nil, fmt.Errorf("create refresh token family key: %w", err)
	}

	record, err := r.client.Get(r.readPolicy, key)

	return refreshTokenFamilyFromRecord(record)
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

func refreshTokenFamilyFromRecord(record *aero.Record) (*RefreshTokenFamily, error) {
	createdAt, err := parseTimeBin(record.Bins[binFamilyCreatedAt], binFamilyCreatedAt)
	if err != nil {
		return nil, err
	}

	familyID, err := uuid.Parse(record.Key.String())

	if err != nil {
		return nil, fmt.Errorf("parse refresh token family ID: %w", err)
	}

	family := &RefreshTokenFamily{
		ID:        familyID,
		CreatedAt: createdAt,
	}

	if value, exists := record.Bins[binFamilyRevokedAt]; exists {
		revokedAt, err := parseTimeBin(value, binFamilyRevokedAt)
		if err != nil {
			return nil, err
		}
		family.RevokedAt = &revokedAt
	}
	return family, nil
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
