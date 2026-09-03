package aerospike

import (
	"fmt"
	"time"

	aero "github.com/aerospike/aerospike-client-go/v8"
	"github.com/fazil-syed/bifrost/internal/config"
)

func newBasePolicy(cfg config.AerospikeConfig) (*aero.BasePolicy, error) {
	socketTimeout, err := time.ParseDuration(cfg.Timeout.Socket)

	if err != nil {
		return nil, fmt.Errorf("parse aerospike socket timeout: %w", err)
	}
	totalTimeout, err := time.ParseDuration(cfg.Timeout.Total)

	if err != nil {
		return nil, fmt.Errorf("parse aerospike total timeout: %w", err)
	}

	sleepBetweenRetries, err := time.ParseDuration(cfg.Timeout.SleepBetweenRetries)

	if err != nil {
		return nil, fmt.Errorf("parse aerospike sleep between retries: %w", err)
	}

	timeoutDelay, err := time.ParseDuration(cfg.Timeout.TimeoutDelay)

	if err != nil {
		return nil, fmt.Errorf("parse aerospike timeout delay: %w", err)
	}

	basePolicy := &aero.BasePolicy{
		SocketTimeout:       socketTimeout,
		TotalTimeout:        totalTimeout,
		MaxRetries:          cfg.Timeout.MaxRetries,
		SleepBetweenRetries: sleepBetweenRetries,
		TimeoutDelay:        timeoutDelay,
	}

	if basePolicy.SocketTimeout > basePolicy.TotalTimeout {
		return nil, fmt.Errorf("aerospike socket timeout %s cannot exceed total timeout %s", basePolicy.SocketTimeout, basePolicy.TotalTimeout)
	}

	return basePolicy, nil
}

func newWritePolicy(cfg config.AerospikeConfig) (*aero.WritePolicy, error) {
	basePolicy, err := newBasePolicy(cfg)
	if err != nil {
		return nil, err
	}
	policy := aero.NewWritePolicy(0, 0)
	policy.BasePolicy = *basePolicy
	return policy, nil
}
