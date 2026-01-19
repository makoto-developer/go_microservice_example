package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetMyReviewsInput represents input for GetMyReviews
type GetMyReviewsInput struct {
	CustomerId uuid.UUID
}

// GetMyReviewsUsecase defines the interface for GetMyReviews
type GetMyReviewsUsecase interface {
	Execute(ctx context.Context, input GetMyReviewsInput) error
}

type get_my_reviewsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetMyReviewsUsecase creates a new instance
func NewGetMyReviewsUsecase() GetMyReviewsUsecase {
	return &get_my_reviewsUsecaseImpl{}
}

// Execute executes GetMyReviews
func (u *get_my_reviewsUsecaseImpl) Execute(ctx context.Context, input GetMyReviewsInput) error {
	// TODO: Implement business logic

	return nil
}
