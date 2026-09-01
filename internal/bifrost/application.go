package bifrost

import (
	"fmt"

	"github.com/fazil-syed/bifrost/internal/authentication"
	"github.com/fazil-syed/bifrost/internal/logger"
	"github.com/fazil-syed/bifrost/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	UserService           user.UserService
	AuthenticationService authentication.AuthenticationService
}

func New(
	db *pgxpool.Pool,
) (*Application, error) {
	userService := user.NewUserService(db)
	authenticationService, err := authentication.NewService(userService)
	if err != nil {
		return nil, fmt.Errorf("initialize authentication service : %w", err)
	}
	return &Application{
		UserService:           userService,
		AuthenticationService: authenticationService,
	}, nil
}

func (a *Application) Start() {
	logger.Info.Println("started bifrost with user service setup")
}
