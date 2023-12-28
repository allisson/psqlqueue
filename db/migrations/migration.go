package migrations

import (
	"embed"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var fs embed.FS

func Migrate(databaseURL string) error {
	slog.Info("migration process started")
	defer slog.Info("migration process finished")

	parsedDatabaseURL := strings.ReplaceAll(databaseURL, "postgresql://", "pgx://")
	parsedDatabaseURL = strings.ReplaceAll(parsedDatabaseURL, "postgres://", "pgx://")

	driver, err := iofs.New(fs, ".")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", driver, parsedDatabaseURL)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
