package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ResolveReviewReportInput represents input for ResolveReviewReport
type ResolveReviewReportInput struct {
	ReportId uuid.UUID
	AdminId uuid.UUID
	Action string
	Note string
}

// ResolveReviewReportUsecase defines the interface for ResolveReviewReport
type ResolveReviewReportUsecase interface {
	Execute(ctx context.Context, input ResolveReviewReportInput) error
}

type resolve_review_reportUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewResolveReviewReportUsecase creates a new instance
func NewResolveReviewReportUsecase() ResolveReviewReportUsecase {
	return &resolve_review_reportUsecaseImpl{}
}

// Execute executes ResolveReviewReport
func (u *resolve_review_reportUsecaseImpl) Execute(ctx context.Context, input ResolveReviewReportInput) error {
	// TODO: Implement business logic

	return nil
}
