package domain

import (
	"context"

	"github.com/google/uuid"
)

// ReviewReportRepository defines repository interface for ReviewReport
type ReviewReportRepository interface {
	// Create creates a new ReviewReport
	Create(ctx context.Context, review_report *ReviewReport) error

	// FindByID finds ReviewReport by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ReviewReport, error)

	// Update updates ReviewReport
	Update(ctx context.Context, review_report *ReviewReport) error

	// Delete deletes ReviewReport by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ReviewReport
	List(ctx context.Context, limit, offset int) ([]*ReviewReport, error)
}
