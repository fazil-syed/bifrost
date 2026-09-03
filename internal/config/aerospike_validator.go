package config

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func validateAerospike(cfg AerospikeConfig) error {
	if len(cfg.Hosts) == 0 {
		return fmt.Errorf("aerospike hosts must contain at least one host")
	}

	for i, host := range cfg.Hosts {
		if err := validateAerospikeHost(i, host); err != nil {
			return err
		}
	}
	if err := validateAerospikeAuthentication(cfg.Authentication); err != nil {
		return err
	}

	if err := valdiateAerospikeTLS(cfg.TLS); err != nil {
		return err
	}

	if err := valdiateAerospikeConnection(cfg.Connection); err != nil {
		return err
	}

	if err := validateAerospikeTimeout(cfg.Timeout); err != nil {
		return err
	}

	return nil
}

func validateAerospikeHost(index int, cfg AerospikeHostConfig) error {
	host := strings.TrimSpace(cfg.Host)

	if host == "" {
		return fmt.Errorf("aerospike hosts[%d].host must not be empty", index)
	}
	if net.ParseIP(host) == nil {
		if strings.ContainsAny(host, " \t\r\n") {
			return fmt.Errorf("aerospike hosts[%d].host contains whitespace", index)
		}
	}

	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("aerospike hosts[%d].port must be between 1 and 65535", index)
	}
	return nil
}

func validateAerospikeAuthentication(cfg AerospikeAuthConfig) error {
	switch cfg.Mode {
	case "none":
		if cfg.Username != "" || cfg.Password != "" {
			return fmt.Errorf("aerospike authentication credentials must be empty when mode is none")
		}

	case "internal":
		if cfg.Username == "" {
			return fmt.Errorf("aerospike authentication username is required for internal authentication")
		}
		if cfg.Password == "" {
			return fmt.Errorf("aerospike authentication password is required for internal authentication")

		}
	default:
		return fmt.Errorf("unsupported aerospike authentication mode: %q", cfg.Mode)
	}
	return nil
}

func valdiateAerospikeTLS(cfg AerospikeTLSConfig) error {
	if !cfg.Enabled {
		return nil
	}

	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("aerospike tls name is required when TLS is enabled")
	}

	if strings.TrimSpace(cfg.CAFile) == "" {
		return fmt.Errorf("aerospike tls ca_file is required when TLS is enabled")
	}
	if (cfg.CertFile == "") != (cfg.KeyFile == "") {
		return fmt.Errorf("aerospike tls cert_file and key_file must be provided together")
	}
	return nil
}

func valdiateAerospikeConnection(cfg AerospikeConnectionConfig) error {
	if cfg.MinConnectionsPerNode < 0 {
		return fmt.Errorf("aerospike min_connections_per_node must be >= 0")
	}

	if cfg.ConnectionQueueSize < 0 {
		return fmt.Errorf("aerospike connection_queue_size must be >= 0")

	}

	if cfg.LimitConnectionsToQueue && cfg.ConnectionQueueSize == 0 {
		return fmt.Errorf("aerospike connection_queue_size must be > 0 when limit_connections_to_queue is enabled")
	}
	return nil
}

func validateAerospikeTimeout(cfg AerospikeTimeoutConfig) error {
	socket, err := parsePositiveDuration(cfg.Socket, "aerospike timeout socket")

	if err != nil {
		return err
	}

	total, err := parsePositiveDuration(cfg.Total, "aerospike timeout total")
	if err != nil {
		return err
	}

	if total < socket {
		return fmt.Errorf("aerospike timeout total must be >= socket timeout")
	}

	if cfg.MaxRetries < 0 {
		return fmt.Errorf("aerospike timeout max_retries must be >= 0")
	}

	if _, err := parseNonNegativeDuration(cfg.SleepBetweenRetries, "aerospike timeout sleep_between_retries"); err != nil {
		return err
	}

	if _, err := parseNonNegativeDuration(cfg.TimeoutDelay, "aerospike timeout timeout_delay"); err != nil {
		return err
	}

	return nil
}

func parsePositiveDuration(value string, field string) (time.Duration, error) {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", field, err)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("%s must be > 0", field)
	}
	return duration, nil
}
func parseNonNegativeDuration(value string, field string) (time.Duration, error) {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", field, err)
	}
	if duration < 0 {
		return 0, fmt.Errorf("%s must be >= 0", field)
	}
	return duration, nil
}
