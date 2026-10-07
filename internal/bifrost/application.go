package bifrost

import (
	"context"
	"fmt"
	"net/http"

	aero "github.com/aerospike/aerospike-client-go/v8"
	"github.com/fazil-syed/bifrost/internal/aerospike"
	authenticationapi "github.com/fazil-syed/bifrost/internal/api/authentication"
	"github.com/fazil-syed/bifrost/internal/api/server"
	"github.com/fazil-syed/bifrost/internal/authentication"
	"github.com/fazil-syed/bifrost/internal/config"
	"github.com/fazil-syed/bifrost/internal/database"
	"github.com/fazil-syed/bifrost/internal/logger"
	"github.com/fazil-syed/bifrost/internal/migrations"
	"github.com/fazil-syed/bifrost/internal/session"
	"github.com/fazil-syed/bifrost/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	db              *pgxpool.Pool
	aerospikeClient *aero.Client
	httpServer      *server.Server
}

func New(
	ctx context.Context,
	cfg config.Config,
) (*Application, error) {

	db, err := database.NewPostgresPool(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	logger.Info.Println("database connection successful")

	if err := migrations.RunGlobal(ctx, db); err != nil {
		db.Close()
		return nil, fmt.Errorf("run global migrations: %w", err)
	}

	logger.Info.Println("global migrations completed")

	if err := migrations.RunAllTenants(ctx, db, cfg.Database); err != nil {
		db.Close()
		return nil, fmt.Errorf("run tenant migrations: %w", err)
	}

	logger.Info.Println("tenant migrations completed")

	aerospikeClient, err := aerospike.New(ctx, cfg.Aerospike)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize aerospike: %w", err)
	}

	logger.Info.Println("aerospike client ready")

	userService := user.NewUserService(db)

	readPolicy, err := aerospike.NewBasePolicy(cfg.Aerospike)

	if err != nil {
		aerospikeClient.Close()
		db.Close()
		return nil, err
	}

	writePolicy, err := aerospike.NewWritePolicy(cfg.Aerospike)
	if err != nil {
		aerospikeClient.Close()
		db.Close()
		return nil, err
	}

	sessionRepository := session.NewAerospikeSessionRepository(aerospikeClient, cfg.Aerospike.Namespace, readPolicy, writePolicy)

	sessionService, err := session.NewSessionService(sessionRepository, cfg.Session.Lifetime)

	if err != nil {
		aerospikeClient.Close()
		db.Close()
		return nil, fmt.Errorf("initialize session service: %w", err)
	}
	authenticationService, err := authentication.NewService(userService, sessionService)
	if err != nil {
		aerospikeClient.Close()
		db.Close()
		return nil, fmt.Errorf("initialize authentication service : %w", err)
	}

	authenticationHandler := authenticationapi.NewHandler(authenticationService)

	mux := http.NewServeMux()

	authenticationapi.RegisterRoutes(mux, authenticationHandler)

	httpServer, err := server.NewServer(cfg.HTTP, mux)

	if err != nil {
		aerospikeClient.Close()
		db.Close()
		return nil, fmt.Errorf("initialize http server: %w", err)
	}

	return &Application{
		db:              db,
		aerospikeClient: aerospikeClient,
		httpServer:      httpServer,
	}, nil
}

func (a *Application) Start(ctx context.Context) error {
	logger.Info.Println("starting bifrost ")

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- a.httpServer.Start()
	}()

	logger.Info.Println("started bifrost server")
	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (a *Application) Shutdown(ctx context.Context) error {
	logger.Info.Println("shutting down bifrost")
	httpErr := a.httpServer.Shutdown(ctx)
	a.aerospikeClient.Close()
	a.db.Close()
	if httpErr != nil {
		return fmt.Errorf("shutdown http server: %w", httpErr)
	}

	logger.Info.Println("bifrost shutdown complete")
	return nil
}
