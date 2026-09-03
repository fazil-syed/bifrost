package aerospike

import (
	"context"
	"fmt"

	aero "github.com/aerospike/aerospike-client-go/v8"
	"github.com/fazil-syed/bifrost/internal/config"
)

func New(ctx context.Context, cfg config.AerospikeConfig) (*aero.Client, error) {
	client, err := newClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize aerospike client %w", err)
	}

	if err := waitForReady(ctx, client); err != nil {
		client.Close()
		return nil, err
	}

	if err := warmUp(client); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}
