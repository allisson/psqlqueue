package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func executeRollback(ctx context.Context, tx pgx.Tx) {
	if err := tx.Rollback(ctx); err != nil {
		slog.Error("database rollback error", "error", err.Error())
	}
}
