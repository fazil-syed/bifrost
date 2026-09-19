package migrations

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/fazil-syed/bifrost/internal/config"
	"github.com/fazil-syed/bifrost/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func RunGlobal(ctx context.Context, pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	migrationFS, err := fs.Sub(Global, "global")

	if err != nil {
		return fmt.Errorf("create global migration filesystem: %w", err)
	}

	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrationFS)
	if err != nil {
		return fmt.Errorf("create global migration provider: %w", err)
	}
	defer provider.Close()

	// Apply the migrations
	_, err = provider.Up(ctx)

	if err != nil {
		return fmt.Errorf("run global migrations: %w", err)
	}
	return nil
}

func RunTenant(ctx context.Context, pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	migrationFS, err := fs.Sub(Tenant, "tenant")

	if err != nil {
		return fmt.Errorf("create tenant migration filesystem: %w", err)
	}

	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrationFS)
	if err != nil {
		return fmt.Errorf("create tenant migration provider: %w", err)
	}
	defer provider.Close()

	// Apply the migrations
	_, err = provider.Up(ctx)

	if err != nil {
		return fmt.Errorf("run tenant migrations: %w", err)
	}
	return nil
}

func RunAllTenants(ctx context.Context, globalPool *pgxpool.Pool, databaseConfig config.DatabaseConfig) error {
	rows, err := globalPool.Query(ctx, `
		SELECT database_name
		FROM tenants
		ORDER BY database_name
	`)

	if err != nil {
		return fmt.Errorf("list tenant databases: %w", err)
	}

	defer rows.Close()
	var databaseNames []string
	for rows.Next() {
		var databaseName string
		if err := rows.Scan(&databaseName); err != nil {
			return fmt.Errorf("scan tenant database: %w", err)
		}
		databaseNames = append(databaseNames, databaseName)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("read tenant databases: %w", err)
	}

	for _, databaseName := range databaseNames {
		tenantDatabaseConfig := databaseConfig
		tenantDatabaseConfig.Name = databaseName

		tenantPool, err := database.NewPostgresPool(ctx, tenantDatabaseConfig)
		if err != nil {
			return fmt.Errorf("connect to tenant database: %q: %w", databaseName, err)
		}

		err = RunTenant(ctx, tenantPool)

		tenantPool.Close()

		if err != nil {
			return fmt.Errorf("run migrations for tenant database %q: %w", databaseName, err)
		}
	}

	return nil
}
