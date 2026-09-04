package session

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	aero "github.com/aerospike/aerospike-client-go/v8"
	"github.com/google/uuid"
)

const sessionSet = "sessions"

const (
	binUserID               = "user_id"
	binAuthenticationMethod = "authentication_method"
	binCreatedAt            = "create_at"
	binExpiresAt            = "expires_at"
	binRevokedAt            = "revoked_at"
)

type AerospikeSessionRepository struct {
	client      *aero.Client
	namespace   string
	readPolicy  *aero.BasePolicy
	writePolicy *aero.WritePolicy
}

func NewAerospikeSession(client *aero.Client, namespace string, readPolicy *aero.BasePolicy, writePolicy *aero.WritePolicy) SessionRepository {
	return &AerospikeSessionRepository{
		client:      client,
		namespace:   namespace,
		readPolicy:  readPolicy,
		writePolicy: writePolicy,
	}

}

func (r *AerospikeSessionRepository) Create(ctx context.Context, session *Session) error {
	key, err := aero.NewKey(r.namespace, sessionSet, session.ID)

	if err != nil {
		return fmt.Errorf("create session key: %w", err)
	}

	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("session has already expired")
	}

	ttlSeconds := uint64(math.Ceil(ttl.Seconds()))

	if ttlSeconds > math.MaxUint32 {
		return fmt.Errorf("session lifetime exceeds aerospike maxinmum TTL")
	}

	policy := *r.writePolicy

	policy.RecordExistsAction = aero.CREATE_ONLY

	policy.Expiration = uint32(ttlSeconds)

	if err = r.client.Put(&policy, key, aero.BinMap{
		binUserID:               session.UserID.String(),
		binAuthenticationMethod: session.AuthenticationMethod,
		binCreatedAt:            session.CreatedAt.UnixNano(),
		binExpiresAt:            session.ExpiresAt.UnixNano(),
	}); err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil

}

func (r *AerospikeSessionRepository) GetByID(ctx context.Context, id string) (*Session, error) {
	key, err := aero.NewKey(r.namespace, sessionSet, id)
	if err != nil {
		return nil, fmt.Errorf("create session key: %w", err)
	}
	record, err := r.client.Get(r.readPolicy, key)
	if err != nil {
		if errors.Is(err, aero.ErrKeyNotFound) {
			return nil, ErrSessionNotFound
		}
	}

	userIDValue, ok := record.Bins[binUserID].(string)

	if !ok {
		return nil, fmt.Errorf("session %q contains invalid user_id", id)
	}

	userID, parseErr := uuid.Parse(userIDValue)
	if parseErr != nil {
		return nil, fmt.Errorf("parse session user_id: %w", parseErr)
	}

	authenticationMethod, ok := record.Bins[binAuthenticationMethod].(string)

	if !ok {
		return nil, fmt.Errorf("session %q contains invalid authentication_method", id)
	}

	createdAtNanos, ok := record.Bins[binCreatedAt].(int64)

	if !ok {
		return nil, fmt.Errorf("session %q contains invalid created_at")
	}
	expiresAtNanos, ok := record.Bins[binExpiresAt].(int64)

	if !ok {
		return nil, fmt.Errorf("session %q contains invalid expires_at")
	}

	session := &Session{
		ID:                   id,
		UserID:               userID,
		AuthenticationMethod: authenticationMethod,
		CreatedAt:            time.Unix(0, createdAtNanos).UTC(),
		ExpiresAt:            time.Unix(0, expiresAtNanos).UTC(),
	}

	if revokedAtNanos, ok := record.Bins[binRevokedAt].(int64); ok {
		revokedAt := time.Unix(0, revokedAtNanos).UTC()
		session.RevokedAt = &revokedAt
	}
	return session, nil
}

func (r *AerospikeSessionRepository) Revoke(ctx context.Context, id string, revokedAt time.Time) error {

	key, err := aero.NewKey(r.namespace, sessionSet, id)

	if err != nil {
		return fmt.Errorf("create session key: %w", err)
	}

	policy := *r.writePolicy

	policy.RecordExistsAction = aero.UPDATE_ONLY
	policy.Expiration = aero.TTLDontUpdate

	if err := r.client.Put(&policy, key, aero.BinMap{
		binRevokedAt: revokedAt.UnixNano(),
	}); err != nil {
		if errors.Is(err, aero.ErrKeyNotFound) {
			return ErrSessionNotFound
		}
		return fmt.Errorf("revoke session: %w", err)
	}
	return nil
}
