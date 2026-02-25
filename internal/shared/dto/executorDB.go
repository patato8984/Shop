package dto

import (
	"context"
	"database/sql"

	"github.com/patato8984/Shop/pkg/ctxkey"
)

type DBQuerier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func Getter(ctx context.Context, defDB *sql.DB) DBQuerier {
	if tx, ok := ctx.Value(ctxkey.TransactionKey).(*sql.Tx); ok {
		return tx
	}
	return defDB
}
