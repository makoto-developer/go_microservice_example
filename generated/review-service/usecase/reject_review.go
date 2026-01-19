package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// RejectReviewInput represents input for RejectReview
type RejectReviewInput struct {
	ReviewId uuid.UUID
	AdminId uuid.UUID
	Reason string
}

// RejectReviewUsecase defines the interface for RejectReview
type RejectReviewUsecase interface {
	Execute(ctx context.Context, input RejectReviewInput) error
}

type reject_reviewUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewRejectReviewUsecase creates a new instance
func NewRejectReviewUsecase() RejectReviewUsecase {
	return &reject_reviewUsecaseImpl{}
}

// Execute executes RejectReview
func (u *reject_reviewUsecaseImpl) Execute(ctx context.Context, input RejectReviewInput) error {
	// TODO: Implement business logic

	return nil
}
