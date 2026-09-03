package bifrost

import (
	"fmt"

	aero "github.com/aerospike/aerospike-client-go/v8"
	"github.com/fazil-syed/bifrost/internal/authentication"
	"github.com/fazil-syed/bifrost/internal/logger"
	"github.com/fazil-syed/bifrost/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	UserService           user.UserService
	AuthenticationService authentication.AuthenticationService
	AerospikeClient       *aero.Client
}

func New(
	db *pgxpool.Pool,
	aerospikeClient *aero.Client,
) (*Application, error) {
	userService := user.NewUserService(db)
	authenticationService, err := authentication.NewService(userService)
	if err != nil {
		return nil, fmt.Errorf("initialize authentication service : %w", err)
	}
	return &Application{
		UserService:           userService,
		AuthenticationService: authenticationService,
		AerospikeClient:       aerospikeClient,
	}, nil
}

func (a *Application) Start() {
	logger.Info.Println("started bifrost with user service setup")
}
