package repos

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
)

func IsNoRowsError(err error) bool {
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return true
	}
	return false
}
