package domain

import (
	"context"

	"github.com/google/uuid"
)

// DashboardMetricsRepository defines repository interface for DashboardMetrics
type DashboardMetricsRepository interface {
	// Create creates a new DashboardMetrics
	Create(ctx context.Context, dashboard_metrics *DashboardMetrics) error

	// FindByID finds DashboardMetrics by ID
	FindByID(ctx context.Context, id uuid.UUID) (*DashboardMetrics, error)

	// Update updates DashboardMetrics
	Update(ctx context.Context, dashboard_metrics *DashboardMetrics) error

	// Delete deletes DashboardMetrics by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all DashboardMetrics
	List(ctx context.Context, limit, offset int) ([]*DashboardMetrics, error)
}
