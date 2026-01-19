package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetReviewsByProductInput represents input for GetReviewsByProduct
type GetReviewsByProductInput struct {
	ProductId uuid.UUID
	SortBy string
	Page int
	PageSize int
}

// GetReviewsByProductUsecase defines the interface for GetReviewsByProduct
type GetReviewsByProductUsecase interface {
	Execute(ctx context.Context, input GetReviewsByProductInput) error
}

type get_reviews_by_productUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetReviewsByProductUsecase creates a new instance
func NewGetReviewsByProductUsecase() GetReviewsByProductUsecase {
	return &get_reviews_by_productUsecaseImpl{}
}

// Execute executes GetReviewsByProduct
func (u *get_reviews_by_productUsecaseImpl) Execute(ctx context.Context, input GetReviewsByProductInput) error {
	// TODO: Implement business logic

	return nil
}
