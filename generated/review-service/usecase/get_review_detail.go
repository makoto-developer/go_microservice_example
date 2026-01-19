package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetReviewDetailInput represents input for GetReviewDetail
type GetReviewDetailInput struct {
	ReviewId uuid.UUID
}

// GetReviewDetailUsecase defines the interface for GetReviewDetail
type GetReviewDetailUsecase interface {
	Execute(ctx context.Context, input GetReviewDetailInput) error
}

type get_review_detailUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetReviewDetailUsecase creates a new instance
func NewGetReviewDetailUsecase() GetReviewDetailUsecase {
	return &get_review_detailUsecaseImpl{}
}

// Execute executes GetReviewDetail
func (u *get_review_detailUsecaseImpl) Execute(ctx context.Context, input GetReviewDetailInput) error {
	// TODO: Implement business logic

	return nil
}
