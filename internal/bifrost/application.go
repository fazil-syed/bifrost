package bifrost

import (
	"github.com/fazil-syed/bifrost/internal/logger"
	"github.com/fazil-syed/bifrost/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	UserService user.UserService
}

func New(
	db *pgxpool.Pool,
) *Application {
	userService := user.NewUserService(db)
	return &Application{
		UserService: userService,
	}
}

func (a *Application) Start() {
	logger.Info.Println("started bifrost with user service setup")
}
