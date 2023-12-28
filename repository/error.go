package repository

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func parseError(err error, notFoundError error, alreadyExistsError error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return notFoundError
	}

	pgerr, ok := err.(*pgconn.PgError)
	if ok {
		switch pgerr.Code {
		case "23505":
			return alreadyExistsError
		}
	}

	return err
}
