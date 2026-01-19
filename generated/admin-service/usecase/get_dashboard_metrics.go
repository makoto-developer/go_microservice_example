package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetDashboardMetricsInput represents input for GetDashboardMetrics
type GetDashboardMetricsInput struct {
	Date date
}

// GetDashboardMetricsUsecase defines the interface for GetDashboardMetrics
type GetDashboardMetricsUsecase interface {
	Execute(ctx context.Context, input GetDashboardMetricsInput) error
}

type get_dashboard_metricsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetDashboardMetricsUsecase creates a new instance
func NewGetDashboardMetricsUsecase() GetDashboardMetricsUsecase {
	return &get_dashboard_metricsUsecaseImpl{}
}

// Execute executes GetDashboardMetrics
func (u *get_dashboard_metricsUsecaseImpl) Execute(ctx context.Context, input GetDashboardMetricsInput) error {
	// TODO: Implement business logic

	return nil
}
