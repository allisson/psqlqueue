package domain

import (
	"log/slog"
	"os"
	"path"

	"github.com/allisson/go-env"
	"github.com/joho/godotenv"
)

func searchup(dir string, filename string) string {
	if dir == "/" || dir == "" {
		return ""
	}

	if _, err := os.Stat(path.Join(dir, filename)); err == nil {
		return path.Join(dir, filename)
	}

	return searchup(path.Dir(dir), filename)
}

func findDotEnv() string {
	directory, err := os.Getwd()
	if err != nil {
		return ""
	}

	filename := ".env"
	return searchup(directory, filename)
}

func loadDotEnv() bool {
	dotenv := findDotEnv()
	if dotenv != "" {
		slog.Info("Found .env", "file", dotenv)
		if err := godotenv.Load(dotenv); err != nil {
			slog.Warn("Can't load .env", "file", dotenv, "error", err)
			return false
		}
		return true
	}
	return false
}

// Config holds all application configuration data.
type Config struct {
	Testing                        bool
	LogLevel                       string
	ServerHost                     string
	ServerPort                     uint
	ServerReadHeaderTimeoutSeconds uint
	MetricsHost                    string
	MetricsPort                    uint
	DatabaseURL                    string
	TestDatabaseURL                string
	DatabaseMinConns               uint
	DatabaseMaxConns               uint
	QueueMaxNumberOfMessages       uint
}

// NewConfig returns a Config with values loaded from environment variables.
func NewConfig() *Config {
	loadDotEnv()

	return &Config{
		Testing:                        env.GetBool("PSQLQUEUE_TESTING", false),
		LogLevel:                       env.GetString("PSQLQUEUE_LOG_LEVEL", "info"),
		ServerHost:                     env.GetString("PSQLQUEUE_SERVER_HOST", "0.0.0.0"),
		ServerPort:                     env.GetUint("PSQLQUEUE_SERVER_PORT", 8000),
		ServerReadHeaderTimeoutSeconds: env.GetUint("PSQLQUEUE_SERVER_READ_HEADER_TIMEOUT_SECONDS", 60),
		MetricsHost:                    env.GetString("PSQLQUEUE_METRICS_HOST", "0.0.0.0"),
		MetricsPort:                    env.GetUint("PSQLQUEUE_METRICS_PORT", 9090),
		DatabaseURL:                    env.GetString("PSQLQUEUE_DATABASE_URL", ""),
		TestDatabaseURL:                env.GetString("PSQLQUEUE_TEST_DATABASE_URL", ""),
		DatabaseMinConns:               env.GetUint("PSQLQUEUE_DATABASE_MIN_CONNS", 0),
		DatabaseMaxConns:               env.GetUint("PSQLQUEUE_DATABASE_MAX_CONNS", 2),
		QueueMaxNumberOfMessages:       env.GetUint("PSQLQUEUE_QUEUE_MAX_NUMBER_OF_MESSAGES", 10),
	}
}
