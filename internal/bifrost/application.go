package bifrost

import (
	"fmt"
	"net/http"

	aero "github.com/aerospike/aerospike-client-go/v8"
	"github.com/fazil-syed/bifrost/internal/aerospike"
	authenticationapi "github.com/fazil-syed/bifrost/internal/api/authentication"
	"github.com/fazil-syed/bifrost/internal/api/server"
	"github.com/fazil-syed/bifrost/internal/authentication"
	"github.com/fazil-syed/bifrost/internal/config"
	"github.com/fazil-syed/bifrost/internal/logger"
	"github.com/fazil-syed/bifrost/internal/session"
	"github.com/fazil-syed/bifrost/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	UserService           user.UserService
	AuthenticationService authentication.AuthenticationService
	SessionService        session.SessionService
	AerospikeClient       *aero.Client
	httpServer            *server.Server
}

func New(
	db *pgxpool.Pool,
	aerospikeClient *aero.Client,
	cfg config.Config,
) (*Application, error) {
	userService := user.NewUserService(db)

	readPolicy, err := aerospike.NewBasePolicy(cfg.Aerospike)

	if err != nil {
		return nil, err
	}

	writePolicy, err := aerospike.NewWritePolicy(cfg.Aerospike)
	if err != nil {
		return nil, err
	}

	sessionRepository := session.NewAerospikeSessionRepository(aerospikeClient, cfg.Aerospike.Namespace, readPolicy, writePolicy)

	sessionService, err := session.NewSessionService(sessionRepository, cfg.Session.Lifetime)

	if err != nil {
		return nil, fmt.Errorf("initialize session service: %w", err)
	}
	authenticationService, err := authentication.NewService(userService, sessionService)
	if err != nil {
		return nil, fmt.Errorf("initialize authentication service : %w", err)
	}

	authenticationHandler := authenticationapi.NewHandler(authenticationService)

	mux := http.NewServeMux()

	authenticationapi.RegisterRoutes(mux, authenticationHandler)

	httpServer, err := server.NewServer(cfg.HTTP, mux)

	if err != nil {
		return nil, fmt.Errorf("initialize http server: %w", err)
	}

	return &Application{
		UserService:           userService,
		AuthenticationService: authenticationService,
		SessionService:        sessionService,
		AerospikeClient:       aerospikeClient,
		httpServer:            httpServer,
	}, nil
}

func (a *Application) Start() {
	logger.Info.Println("starting bifrost ")

	if err := a.httpServer.Start(); err != nil {
		logger.Error.Fatalf("http server stopped unexpectedly: %v", err)
	}

}
