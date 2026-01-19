package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpdateReviewInput represents input for UpdateReview
type UpdateReviewInput struct {
	ReviewId uuid.UUID
	CustomerId uuid.UUID
	Rating int
	Title string
	Content string
	ImageUrls []string
}

// UpdateReviewUsecase defines the interface for UpdateReview
type UpdateReviewUsecase interface {
	Execute(ctx context.Context, input UpdateReviewInput) error
}

type update_reviewUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUpdateReviewUsecase creates a new instance
func NewUpdateReviewUsecase() UpdateReviewUsecase {
	return &update_reviewUsecaseImpl{}
}

// Execute executes UpdateReview
func (u *update_reviewUsecaseImpl) Execute(ctx context.Context, input UpdateReviewInput) error {
	// TODO: Implement business logic

	return nil
}
