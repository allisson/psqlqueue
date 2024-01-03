package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/allisson/psqlqueue/domain"
)

// HealthCheck is an implementation of domain.HealthCheckRepository.
type HealthCheck struct {
	pool *pgxpool.Pool
}

func (h *HealthCheck) Check(ctx context.Context) (*domain.HealthCheck, error) {
	result := 0
	check := &domain.HealthCheck{}

	if err := h.pool.QueryRow(ctx, "SELECT 1+1").Scan(&result); err != nil {
		return check, err
	}

	check.Success = result == 2
	return check, nil
}

// NewHealthCheck returns an implementation of domain.HealthCheckRepository.
func NewHealthCheck(pool *pgxpool.Pool) *HealthCheck {
	return &HealthCheck{pool: pool}
}
