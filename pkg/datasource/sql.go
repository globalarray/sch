package datasource

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"benzo/pkg/utils"

	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"
)

type (
	Conn interface {
		BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
		PingContext(ctx context.Context) (err error)
		io.Closer
		ConnTx
	}

	ConnTx interface {
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
		QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
		QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	}

	Exec interface {
		Scan(rowsAffected, lastInsertID *int64) (err error)
	}

	Query interface {
		Scan(row func(i int) utils.Array) (err error)
	}

	exec struct {
		sqlResult sql.Result
		err       error
	}

	query struct {
		sqlRows *sqlx.Rows
		err     error
	}

	DataSource struct{}
)

var (
	_   Conn   = (*sqlx.Conn)(nil)
	_   Conn   = (*sqlx.DB)(nil)
	_   ConnTx = (*sqlx.Tx)(nil)
	log        = slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelError,
		TimeFormat: time.Kitchen,
	}))
)

// datasource errors
var (
	ErrNoColumnReturned = errors.New("no columns returned")
	ErrDataNotFound     = errors.New("data not found")
	ErrInvalidArguments = errors.New("invalid arguments for scan")
)

func (x exec) Scan(rowsAffected, lastInsertID *int64) error {
	if x.err != nil {
		log.Error("[database:exec]error not nil", slog.Any("error", x.err))

		return x.err
	}

	if x.sqlResult == nil {
		log.Error("[database:exec]rows is nil", slog.Any("error", sql.ErrNoRows))

		return ErrDataNotFound
	}

	if rowsAffected != nil {
		n, err := x.sqlResult.RowsAffected()
		if err != nil {
			log.Error("[database:exec]scan rows affected error", slog.Any("error", err))

			return err
		}
		if n < 1 {
			log.Error("[database:exec]scan rows affected error", slog.Any("error", ErrDataNotFound))

			return ErrDataNotFound
		}
		*rowsAffected = int64(n)
	}

	if lastInsertID != nil {
		n, err := x.sqlResult.LastInsertId()
		if err != nil {
			log.Error("[database:exec]last inserted id error", slog.Any("error", err))
		} else {
			*lastInsertID = int64(n)
		}
	}

	return nil
}

func (x query) Scan(row func(i int) utils.Array) error {
	if x.err != nil {
		log.Error("[database:query]error not nil", slog.Any("error", x.err))

		return x.err
	}

	if x.sqlRows == nil {
		log.Error("[database:query]rows is nil", slog.Any("error", sql.ErrNoRows))

		return ErrDataNotFound
	}

	if err := x.sqlRows.Err(); err != nil {
		return err
	}

	defer x.sqlRows.Close()

	columns, err := x.sqlRows.Columns()
	if err != nil {
		log.Error("[database:query]columns", slog.Any("error", err))

		return err
	}

	if len(columns) < 1 {
		log.Error("[database:query]count columns length", slog.Any("error", ErrNoColumnReturned))

		return ErrNoColumnReturned
	}

	var idx int = 0
	for x.sqlRows.Next() {
		if x.sqlRows.Err() != nil {
			log.Error("[database:query]error to scan sql rows", slog.Any("error", x.sqlRows.Err()))

			return x.sqlRows.Err()
		}

		if row(idx) == nil {
			break
		}

		if len(row(idx)) < 1 {
			continue
		}

		if len(row(idx)) != len(columns) {
			err := fmt.Errorf("%w: [%d] columns on [%d] destinations", ErrInvalidArguments, len(columns), len(row(idx)))
			log.Error("[database:query]error invalid args to scan", slog.Any("error", err))

			return err
		}

		if err = x.sqlRows.Scan(row(idx)...); err != nil {
			log.Error("[database:query] failed to scan row", slog.Any("error", err))

			return err
		}

		idx++
	}

	return err
}

func (DataSource) ExecSQL(sqlResult sql.Result, err error) exec {
	return exec{sqlResult, err}
}

func (DataSource) QuerySQL(sqlRows *sqlx.Rows, err error) Query {
	return query{sqlRows, err}
}
