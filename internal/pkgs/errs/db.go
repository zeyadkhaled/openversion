package errs

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

const pgErrUniqueConstraint = "23505"

type Op uint8

const (
	OpOther Op = iota
	OpINSERT
	OpUPDATE
)

func Postgre(err error, op Op) error {
	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == pgErrUniqueConstraint {
		k := KindConflict
		if op == OpINSERT {
			k = KindDuplicate
		}

		return E{
			Kind:    k,
			Wrapped: err,
		}

	}

	if err == pgx.ErrNoRows {
		return E{Kind: KindNotFound}
	}

	return err
}
