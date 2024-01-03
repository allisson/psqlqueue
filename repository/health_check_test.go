package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/allisson/psqlqueue/domain"
)

func TestHealthCheck(t *testing.T) {
	cfg := domain.NewConfig()
	ctx := context.Background()
	pool, _ := pgxpool.New(ctx, cfg.TestDatabaseURL)
	defer pool.Close()

	t.Run("Check", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		healthCheckRepo := NewHealthCheck(pool)

		healthCheck, err := healthCheckRepo.Check(ctx)
		assert.Nil(t, err)
		assert.True(t, healthCheck.Success)
	})
}
