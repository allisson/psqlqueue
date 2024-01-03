package service

import (
	"context"

	"github.com/allisson/psqlqueue/domain"
)

// HealthCheck is an implementation of domain.HealthCheckService.
type HealthCheck struct {
	healthCheckRepository domain.HealthCheckRepository
}

func (h *HealthCheck) Check(ctx context.Context) (*domain.HealthCheck, error) {
	return h.healthCheckRepository.Check(ctx)
}

// NewHealthCheck returns an implementation of domain.HealthCheckService.
func NewHealthCheck(healthCheckRepository domain.HealthCheckRepository) *HealthCheck {
	return &HealthCheck{healthCheckRepository: healthCheckRepository}
}
