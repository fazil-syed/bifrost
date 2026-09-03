package aerospike

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/fazil-syed/bifrost/internal/config"

	aero "github.com/aerospike/aerospike-client-go/v8"
)

func newClientPolicy(cfg config.AerospikeConfig) (*aero.ClientPolicy, error) {
	policy := aero.NewClientPolicy()

	// Connection pool configuraiton
	policy.MinConnectionsPerNode = cfg.Connection.MinConnectionsPerNode
	policy.ConnectionQueueSize = cfg.Connection.ConnectionQueueSize
	policy.LimitConnectionsToQueueSize = cfg.Connection.LimitConnectionsToQueue

	// Authentication

	switch cfg.Authentication.Mode {
	case "none":
		// Leave User and Password empty
		// ClientPolicy defaults AuthMode to Internal

	case "internal":
		policy.AuthMode = aero.AuthModeInternal
		policy.User = cfg.Authentication.Username
		policy.Password = cfg.Authentication.Password

	default:
		return nil, fmt.Errorf("unsupported aerospike authentication mode %q", cfg.Authentication.Mode)
	}

	tlsConfig, err := newTlSConfig(cfg.TLS)
	if err != nil {
		return nil, err
	}

	policy.TlsConfig = tlsConfig

	return policy, nil
}

func newTlSConfig(cfg config.AerospikeTLSConfig) (*tls.Config, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	if cfg.Name == "" {
		return nil, fmt.Errorf("aerospike tls name is required when TLS is enabled")
	}
	if cfg.CAFile == "" {
		return nil, fmt.Errorf("aerospike tls ca_file is required when TLS is enabled")
	}

	//Load the CA certificate used to verify the Aerospike server
	caPEM, err := os.ReadFile(cfg.CAFile)

	if err != nil {
		return nil, fmt.Errorf("read aerospike tls CA file %q: %w", cfg.CAFile, err)
	}

	rootCAs := x509.NewCertPool()

	if ok := rootCAs.AppendCertsFromPEM(caPEM); !ok {
		return nil, fmt.Errorf("parse aerospike tls CA certificate %q: no certificates found", cfg.CAFile)
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: cfg.Name,
		RootCAs:    rootCAs,
	}

	if (cfg.CertFile == "") != (cfg.KeyFile == "") {
		return nil, fmt.Errorf("aerospike tls cert_file and key_file must be provided together")
	}
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		clientCert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("load aerospike tls client certificate and key: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}

	return nil, nil
}
