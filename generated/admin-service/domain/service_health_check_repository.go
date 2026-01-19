package domain

import (
	"context"

	"github.com/google/uuid"
)

// ServiceHealthCheckRepository defines repository interface for ServiceHealthCheck
type ServiceHealthCheckRepository interface {
	// Create creates a new ServiceHealthCheck
	Create(ctx context.Context, service_health_check *ServiceHealthCheck) error

	// FindByID finds ServiceHealthCheck by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ServiceHealthCheck, error)

	// Update updates ServiceHealthCheck
	Update(ctx context.Context, service_health_check *ServiceHealthCheck) error

	// Delete deletes ServiceHealthCheck by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ServiceHealthCheck
	List(ctx context.Context, limit, offset int) ([]*ServiceHealthCheck, error)
}
