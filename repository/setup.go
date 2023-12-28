package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/allisson/psqlqueue/domain"
)

// SetupDatabaseConnection returns a configured pgx pool.
func SetupDatabaseConnection(ctx context.Context, cfg *domain.Config) (*pgxpool.Pool, error) {
	databaseURL := cfg.DatabaseURL
	if cfg.Testing {
		databaseURL = cfg.TestDatabaseURL
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		slog.Error("database config error", "error", err.Error())
		return nil, err
	}
	config.MinConns = int32(cfg.DatabaseMinConns)
	config.MaxConns = int32(cfg.DatabaseMaxConns)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error("database pool error", "error", err.Error())
		return nil, err
	}

	return pool, nil
}
