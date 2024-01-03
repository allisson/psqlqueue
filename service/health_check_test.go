package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/allisson/psqlqueue/domain"
	"github.com/allisson/psqlqueue/mocks"
)

func TestHealthCheck(t *testing.T) {
	ctx := context.Background()

	t.Run("Check", func(t *testing.T) {
		healthCheckRepository := mocks.NewHealthCheckRepository(t)
		healthCheckService := NewHealthCheck(healthCheckRepository)

		healthCheckRepository.On("Check", ctx).Return(&domain.HealthCheck{Success: true}, nil)

		healthCheck, err := healthCheckService.Check(ctx)
		assert.Nil(t, err)
		assert.True(t, healthCheck.Success)
	})
}
