package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ReportReviewInput represents input for ReportReview
type ReportReviewInput struct {
	ReviewId uuid.UUID
	ReporterId uuid.UUID
	Reason ReportReason
	Description string
}

// ReportReviewUsecase defines the interface for ReportReview
type ReportReviewUsecase interface {
	Execute(ctx context.Context, input ReportReviewInput) error
}

type report_reviewUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewReportReviewUsecase creates a new instance
func NewReportReviewUsecase() ReportReviewUsecase {
	return &report_reviewUsecaseImpl{}
}

// Execute executes ReportReview
func (u *report_reviewUsecaseImpl) Execute(ctx context.Context, input ReportReviewInput) error {
	// TODO: Implement business logic

	return nil
}
