package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fazil-syed/bifrost/internal/config"
	"github.com/fazil-syed/bifrost/internal/logger"
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
func (s *Server) Start() error {
	logger.Info.Println("starting bifrost server")
	err := s.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)

	defer cancel()

	return s.server.Shutdown(shutdownCtx)
}
