package main

import (
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/allisson/psqlqueue/db/migrations"
	"github.com/allisson/psqlqueue/domain"
)

func main() {
	cfg := domain.NewConfig()
	logger := domain.NewLogger(cfg.LogLevel)
	slog.SetDefault(logger)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "migrate",
				Aliases: []string{"m"},
				Usage:   "run database migrate",
				Action: func(c *cli.Context) error {
					if cfg.Testing {
						return migrations.Migrate(cfg.TestDatabaseURL)
					}
					return migrations.Migrate(cfg.DatabaseURL)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("cli app failed", "error", err)
		os.Exit(1)
	}
}
