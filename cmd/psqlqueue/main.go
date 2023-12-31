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
					topicRepository := repository.NewTopic(pool)
					subscriptionRepository := repository.NewSubscription(pool)
					healthCheckRepository := repository.NewHealthCheck(pool)

					// services
					queueService := service.NewQueue(queueRepository)
					messageService := service.NewMessage(messageRepository, queueRepository)
					topicService := service.NewTopic(topicRepository, subscriptionRepository, queueRepository, messageRepository)
					subscriptionService := service.NewSubscription(subscriptionRepository)
					healthCheckService := service.NewHealthCheck(healthCheckRepository)

					// http handlers
					queueHandler := http.NewQueueHandler(queueService)
					messageHandler := http.NewMessageHandler(messageService)
					topicHandler := http.NewTopicHandler(topicService)
					subscriptionHandler := http.NewSubscriptionHandler(subscriptionService)
					healthCheckHandler := http.NewHealthCheckHandler(healthCheckService)

					// run http server
					http.RunServer(c.Context, cfg, http.SetupRouter(logger, queueHandler, messageHandler, topicHandler, subscriptionHandler, healthCheckHandler))

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
