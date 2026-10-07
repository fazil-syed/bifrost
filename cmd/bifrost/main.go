package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fazil-syed/bifrost/internal/bifrost"
	"github.com/fazil-syed/bifrost/internal/config"
	"github.com/fazil-syed/bifrost/internal/logger"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer stop()

	app, err := bifrost.New(ctx, *cfg)

	if err != nil {
		logger.Error.Fatalf("failed to initialize bifrst : %v", err)
	}

	if err := app.Start(ctx); err != nil {
		logger.Error.Fatal(err)
	}

	if err := app.Shutdown(context.Background()); err != nil {
		logger.Error.Fatal(err)
	}

}
