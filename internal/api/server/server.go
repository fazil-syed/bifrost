package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fazil-syed/bifrost/internal/config"
)

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

func NewServer(cfg config.HTTPConfig, handler http.Handler) (*Server, error) {
	if handler == nil {
		return nil, fmt.Errorf("http handler is required")
	}

	return &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		shutdownTimeout: cfg.ShutdownTimeout,
	}, nil
}
