package migrations

import (
	"context"
	"fmt"
	"io/fs"

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
		return fmt.Errorf("create migration provider: %w", err)
	}
	defer provider.Close()

	// Apply the migrations
	_, err = provider.Up(ctx)

	if err != nil {
		return fmt.Errorf("run global migrations: %w", err)
	}
	return nil
}
