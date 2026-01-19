package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetReviewReportsInput represents input for GetReviewReports
type GetReviewReportsInput struct {
	AdminId uuid.UUID
	Status ReportStatus
	Page int
	PageSize int
}

// GetReviewReportsUsecase defines the interface for GetReviewReports
type GetReviewReportsUsecase interface {
	Execute(ctx context.Context, input GetReviewReportsInput) error
}

type get_review_reportsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetReviewReportsUsecase creates a new instance
func NewGetReviewReportsUsecase() GetReviewReportsUsecase {
	return &get_review_reportsUsecaseImpl{}
}

// Execute executes GetReviewReports
func (u *get_review_reportsUsecaseImpl) Execute(ctx context.Context, input GetReviewReportsInput) error {
	// TODO: Implement business logic

	return nil
}
