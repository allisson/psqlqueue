package main

import (
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/allisson/psqlqueue/db/migrations"
	"github.com/allisson/psqlqueue/domain"
	"github.com/allisson/psqlqueue/http"
	"github.com/allisson/psqlqueue/repository"
	"github.com/allisson/psqlqueue/service"
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
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "run http server",
				Action: func(c *cli.Context) error {
					pool, err := repository.SetupDatabaseConnection(c.Context, cfg)
					if err != nil {
						return err
					}
					defer pool.Close()

					// repositories
					queueRepository := repository.NewQueue(pool)
					messageRepository := repository.NewMessage(pool)

					// services
					queueService := service.NewQueue(queueRepository)
					messageService := service.NewMessage(messageRepository, queueRepository)

					// http handlers
					queueHandler := http.NewQueueHandler(queueService)
					messageHandler := http.NewMessageHandler(messageService)

					// run http server
					http.RunServer(c.Context, cfg, http.SetupRouter(logger, queueHandler, messageHandler))

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("cli app failed", "error", err)
		os.Exit(1)
	}
}
