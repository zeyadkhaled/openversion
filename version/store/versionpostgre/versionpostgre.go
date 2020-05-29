package versionpostgre

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/api/global"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/errs"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version"
)

type Store struct {
	pool   *pgxpool.Pool
	logger zerolog.Logger
}

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func get(ctx context.Context, q querier, id string) (a version.Application, err error) {
	ctx, span := global.Tracer("service").Start(ctx, "store.postgre.get")
	defer span.End()

	const query = `SELECT "id", "min_version", "package","created_at", "updated_at"
		 FROM backend.versions WHERE "id" = $1`
	row := q.QueryRow(ctx, query, id)
	err = row.Scan(&a.ID, &a.MinVersion, &a.Package, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return version.Application{}, errs.Postgre(err, errs.OpOther)
	}
	return a, nil
}

func New(ctx context.Context, connStr string, logger zerolog.Logger) (*Store, error) {
	cc, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}
	cc.ConnConfig.Logger = zerologadapter.NewLogger(logger)

	pgxPool, err := pgxpool.ConnectConfig(ctx, cc)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}

	store := &Store{
		pool:   pgxPool,
		logger: logger,
	}

	return store, nil
}

func (store *Store) Upsert(ctx context.Context, a version.Application) error {
	ctx, span := global.Tracer("service").Start(ctx, "store.postgre.Upsert")
	defer span.End()

	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			store.logger.Error().Err(err).Msg("failed to rollback application version update")
		}
	}()

	_, err = get(ctx, tx, a.ID)
	if errs.Is(err, errs.KindNotFound) {
		_, err = tx.Exec(ctx, `INSERT INTO backend.versions (
			"id", "min_version", "package",
			"created_at", "updated_at")
			VALUES ($1, $2, $3, $4, $5)`,
			a.ID, a.MinVersion, a.Package, a.CreatedAt, a.UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to insert application version: %w", errs.Postgre(err, errs.OpINSERT))
		}
	} else {
		tag, err := tx.Exec(ctx, `UPDATE backend.versions SET (
			"id", "min_version", "package","created_at", "updated_at") =
			($1, $2, $3, $4, $5) WHERE "id" = $6`,
			a.ID, a.MinVersion, a.Package, a.CreatedAt, a.UpdatedAt, a.ID)
		if err != nil {
			return fmt.Errorf("failed to update application: %w", errs.Postgre(err, errs.OpUPDATE))
		}
		if tag.RowsAffected() == 0 {
			return errs.E{Kind: errs.KindConflict}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit version.Insert: %w", errs.Postgre(err, errs.OpOther))
	}

	return nil
}

func (store *Store) Get(ctx context.Context, id string) (a version.Application, err error) {
	ctx, span := global.Tracer("service").Start(ctx, "store.postgre.Get")
	defer span.End()

	return get(ctx, store.pool, id)
}

func (store *Store) List(ctx context.Context, limit int) (apps []version.Application, err error) {
	ctx, span := global.Tracer("service").Start(ctx, "store.postgre.List")
	defer span.End()

	query := `SELECT "id", "min_version", "package",
			 "created_at", "updated_at"
			 FROM backend.versions`

	rows, err := store.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.Postgre(err, errs.OpOther)
	}
	defer rows.Close()

	for rows.Next() {
		var a version.Application

		err := rows.Scan(
			&a.ID, &a.MinVersion, &a.Package, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, errs.Postgre(err, errs.OpOther)
		}

		apps = append(apps, a)
	}

	return apps, nil
}
