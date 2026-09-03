package aerospike

import (
	"fmt"

	aero "github.com/aerospike/aerospike-client-go/v8"
)

func warmUp(client *aero.Client) error {
	if _, err := client.WarmUp(0); err != nil {
		return fmt.Errorf("warm up aerospike connection pool: %w", err)
	}
	return nil
}
