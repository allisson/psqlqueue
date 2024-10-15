package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	sloggin "github.com/samber/slog-gin"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/allisson/psqlqueue/docs"
	"github.com/allisson/psqlqueue/domain"
)

var prometheusMiddleware = middleware.New(middleware.Config{
	Recorder: metrics.NewRecorder(metrics.Config{}),
})

func SetupRouter(logger *slog.Logger, queueHandler *QueueHandler, messageHandler *MessageHandler, topicHandler *TopicHandler, subscriptionHandler *SubscriptionHandler, healthCheckHandler *HealthCheckHandler) *gin.Engine {
	// router setup
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(sloggin.New(logger), gin.Recovery())
	router.Use(ginmiddleware.Handler("", prometheusMiddleware))

	// swagger config
	docs.SwaggerInfo.Title = "PSQL Queue API"
	docs.SwaggerInfo.BasePath = "/v1"

	// v1 group setup
	v1 := router.Group("/v1")

	// swagger
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// queue handler
	v1.POST("/queues", queueHandler.Create)
	v1.PUT("/queues/:queue_id", queueHandler.Update)
	v1.GET("/queues/:queue_id", queueHandler.Get)
	v1.GET("/queues", queueHandler.List)
	v1.DELETE("/queues/:queue_id", queueHandler.Delete)
	v1.GET("/queues/:queue_id/stats", queueHandler.Stats)
	v1.PUT("/queues/:queue_id/purge", queueHandler.Purge)
	v1.PUT("/queues/:queue_id/cleanup", queueHandler.Cleanup)

	// message handler
	v1.POST("/queues/:queue_id/messages", messageHandler.Create)
	v1.GET("/queues/:queue_id/messages", messageHandler.List)
	v1.PUT("/queues/:queue_id/messages/:message_id/ack", messageHandler.Ack)
	v1.PUT("/queues/:queue_id/messages/:message_id/nack", messageHandler.Nack)

	// topic handler
	v1.POST("/topics", topicHandler.Create)
	v1.GET("/topics/:topic_id", topicHandler.Get)
	v1.GET("/topics", topicHandler.List)
	v1.DELETE("/topics/:topic_id", topicHandler.Delete)
	v1.POST("/topics/:topic_id/messages", topicHandler.CreateMessage)

	// subscription handler
	v1.POST("/subscriptions", subscriptionHandler.Create)
	v1.GET("/subscriptions/:subscription_id", subscriptionHandler.Get)
	v1.GET("/subscriptions", subscriptionHandler.List)
	v1.DELETE("/subscriptions/:subscription_id", subscriptionHandler.Delete)

	// health check handler
	v1.GET("/healthz", healthCheckHandler.Check)

	return router
}

// RunServer runs an HTTP server based on config and router.
func RunServer(ctx context.Context, cfg *domain.Config, router *gin.Engine) {

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	metricsSrv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.MetricsHost, cfg.MetricsPort),
		Handler:           promhttp.Handler(),
		ReadHeaderTimeout: time.Second * time.Duration(cfg.ServerReadHeaderTimeoutSeconds),
	}
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort),
		Handler:           router,
		ReadHeaderTimeout: time.Second * time.Duration(cfg.ServerReadHeaderTimeoutSeconds),
	}

	// Initializing the metrics server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		slog.Info("metrics server starting", "host", cfg.MetricsHost, "port", cfg.MetricsPort)

		if err := metricsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("metrics server listen error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Initializing the http server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		slog.Info("http server starting", "host", cfg.ServerHost, "port", cfg.ServerPort)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server listen error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	slog.Info("server shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := metricsSrv.Shutdown(ctx); err != nil {
		slog.Error("metrics server forced to shutdown: ", slog.Any("error", err))
		os.Exit(1)
	}
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("http server forced to shutdown: ", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("server exiting")
}
