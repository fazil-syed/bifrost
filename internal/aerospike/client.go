package aerospike

import (
	"fmt"

	aero "github.com/aerospike/aerospike-client-go/v8"
	"github.com/fazil-syed/bifrost/internal/config"
)

func newClient(cfg config.AerospikeConfig) (*aero.Client, error) {
	policy, err := newClientPolicy(cfg)
	if err != nil {
		return nil, fmt.Errorf("create aerospike client policy: %w", err)
	}

	hosts := make([]*aero.Host, 0, len(cfg.Hosts))

	for _, host := range cfg.Hosts {
		hosts = append(hosts, aero.NewHost(host.Host, host.Port))
	}

	client, err := aero.NewClientWithPolicyAndHost(policy, hosts...)
	if err != nil {
		return nil, fmt.Errorf("create aerospike client: %w", err)
	}
	return client, nil
}
