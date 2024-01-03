package domain

import "context"

// HealthCheck entity.
type HealthCheck struct {
	Success bool `json:"success"`
}

// HealthCheckRepository is the repository interface for the HealthCheck entity.
type HealthCheckRepository interface {
	Check(ctx context.Context) (*HealthCheck, error)
}

// HealthCheckService is the service interface for the HealthCheck entity.
type HealthCheckService interface {
	Check(ctx context.Context) (*HealthCheck, error)
}
