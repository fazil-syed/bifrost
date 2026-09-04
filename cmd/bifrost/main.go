package main

import (
	"context"
	"flag"
	"log"

	"github.com/fazil-syed/bifrost/internal/aerospike"
	"github.com/fazil-syed/bifrost/internal/bifrost"
	"github.com/fazil-syed/bifrost/internal/config"
	"github.com/fazil-syed/bifrost/internal/database"
	"github.com/fazil-syed/bifrost/internal/logger"
	"github.com/fazil-syed/bifrost/internal/migrations"
)

func main() {
	log.Printf("App initiating")

	configPath := flag.String("config", "config.yaml", "set the config yaml file path")

	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	if err := config.Validate(cfg); err != nil {
		panic(err)
	}

	logger.Init(cfg.Bifrost.Name, cfg.Logging.Level)

	logger.Info.Printf("Bifrost logger initialized")

	ctx := context.Background()
	db, err := database.NewPostgresPool(ctx, cfg.Database)
	if err != nil {
		logger.Error.Fatalf("initialize database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		logger.Error.Fatalf("ping database: %v", err)
	}

	logger.Info.Println("database connection successful")

	if err := migrations.RunGlobal(ctx, db); err != nil {
		logger.Error.Fatalf("run global migrations: %v", err)
	}

	logger.Info.Println("global migrations completed")

	aerospikeClient, err := aerospike.New(ctx, cfg.Aerospike)
	if err != nil {
		logger.Error.Fatalf("initialize aerospike: %w", err)
	}
	defer aerospikeClient.Close()

	logger.Info.Println("aerospike client ready")

	app, err := bifrost.New(db, aerospikeClient, *cfg)

	if err != nil {
		logger.Error.Fatalf("failed to initialize bifrst : %v", err)
	}

	app.Start()

}
