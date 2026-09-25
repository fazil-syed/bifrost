package database

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsUniqueViolation(err error) (*pgconn.PgError, bool) {
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return nil, false
	}

	return pgErr, pgErr.Code == pgerrcode.UniqueViolation
}
