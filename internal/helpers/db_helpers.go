package helpers

import (
	"database/sql"
	"github.com/jackc/pgx"
	"github.com/yletamitlu/tech-db/internal/consts"
)

const (
	PathItemLen        = 8
	NullPathItem       = "00000000"
	MaxNesting         = 5
	PathItemsSeparator = "."

	DeadlockErrorCode = "40P01"
)

func PgxErrToCustom(err error) error {
	switch err {
	case pgx.ErrNoRows:
		return consts.ErrNotFound
	case sql.ErrNoRows:
		return consts.ErrNotFound
	default:
		return err
	}
}
