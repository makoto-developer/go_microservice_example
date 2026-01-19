package domain

import (
	"context"

	"github.com/google/uuid"
)

// SalesReportRepository defines repository interface for SalesReport
type SalesReportRepository interface {
	// Create creates a new SalesReport
	Create(ctx context.Context, sales_report *SalesReport) error

	// FindByID finds SalesReport by ID
	FindByID(ctx context.Context, id uuid.UUID) (*SalesReport, error)

	// Update updates SalesReport
	Update(ctx context.Context, sales_report *SalesReport) error

	// Delete deletes SalesReport by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all SalesReport
	List(ctx context.Context, limit, offset int) ([]*SalesReport, error)
}
