package aerospike

import (
	"context"
	"fmt"
	"time"

	aero "github.com/aerospike/aerospike-client-go/v8"
)

const readinessPollInterval = 100 * time.Millisecond

func waitForReady(ctx context.Context, client *aero.Client) error {
	ticker := time.NewTicker(readinessPollInterval)

	defer ticker.Stop()
	for {
		if client.IsConnected() {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("aerospike client readiness failed: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}
