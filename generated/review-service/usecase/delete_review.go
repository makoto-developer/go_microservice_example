package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DeleteReviewInput represents input for DeleteReview
type DeleteReviewInput struct {
	ReviewId uuid.UUID
	CustomerId uuid.UUID
}

// DeleteReviewUsecase defines the interface for DeleteReview
type DeleteReviewUsecase interface {
	Execute(ctx context.Context, input DeleteReviewInput) error
}

type delete_reviewUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewDeleteReviewUsecase creates a new instance
func NewDeleteReviewUsecase() DeleteReviewUsecase {
	return &delete_reviewUsecaseImpl{}
}

// Execute executes DeleteReview
func (u *delete_reviewUsecaseImpl) Execute(ctx context.Context, input DeleteReviewInput) error {
	// TODO: Implement business logic

	return nil
}
