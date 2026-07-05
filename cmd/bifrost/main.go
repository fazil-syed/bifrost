package main

import (
	"flag"
	"log"

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

}
