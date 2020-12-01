package helpers

import (
	"database/sql"
	"github.com/jackc/pgx"
	"github.com/yletamitlu/tech-db/internal/consts"
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
