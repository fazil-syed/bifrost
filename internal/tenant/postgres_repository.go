package tenant

import "github.com/jackc/pgx/v5"

type postgresTenantRepository struct {
	tx pgx.Tx
}
