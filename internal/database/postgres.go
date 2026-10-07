package database

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/fazil-syed/bifrost/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func buildVerifiedTLSConfig(cfg config.SSLConfig) (*tls.Config, error) {
	caData, err := os.ReadFile(cfg.RootCert)
	if err != nil {
		return nil, fmt.Errorf("read database root certificate: %w", err)
	}

	rootCAs := x509.NewCertPool()

	if !rootCAs.AppendCertsFromPEM(caData) {
		return nil, fmt.Errorf("failed to parse database root certificate")
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    rootCAs,
	}

	if cfg.Cert != "" || cfg.Key != "" {
		cert, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
		if err != nil {
			return nil, fmt.Errorf("load database client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	return tlsConfig, nil
}

func configureTLS(connConfig *pgx.ConnConfig, cfg config.SSLConfig) error {
	switch cfg.Mode {
	case "disable":
		connConfig.TLSConfig = nil
	case "require":
		connConfig.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	case "verify-ca", "verify-full":
		tlsConfig, err := buildVerifiedTLSConfig(cfg)
		if err != nil {
			return err
		}
		if cfg.Mode == "verify-full" {
			tlsConfig.ServerName = connConfig.Host
		}
		connConfig.TLSConfig = tlsConfig

	default:
		return fmt.Errorf("unsupported database ssl mode: %q", cfg.Mode)
	}
	return nil
}

func NewPostgresPoolConfig(cfg config.DatabaseConfig) (*pgxpool.Config, error) {
	poolConfig, err := pgxpool.ParseConfig("")

	if err != nil {
		return nil, fmt.Errorf("parse postgres pool config: %w", err)
	}

	connConfig := poolConfig.ConnConfig

	connConfig.Host = cfg.Host
	connConfig.Port = uint16(cfg.Port)
	connConfig.Database = cfg.Name
	connConfig.User = cfg.User
	connConfig.Password = cfg.Password

	if err := configureTLS(connConfig, cfg.SSL); err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(cfg.Pool.MaxConns)
	poolConfig.MinConns = int32(cfg.Pool.MinConns)
	poolConfig.MaxConnLifetime = cfg.Pool.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.Pool.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.Pool.HealthCheckPeriod

	connConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
	return poolConfig, nil
}

func NewPostgresPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolConfig, err := NewPostgresPoolConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("configure postgres : %w", err)
	}
	return pgxpool.NewWithConfig(ctx, poolConfig)
}
