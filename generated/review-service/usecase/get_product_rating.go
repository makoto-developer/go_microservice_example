package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetProductRatingInput represents input for GetProductRating
type GetProductRatingInput struct {
	ProductId uuid.UUID
}

// GetProductRatingUsecase defines the interface for GetProductRating
type GetProductRatingUsecase interface {
	Execute(ctx context.Context, input GetProductRatingInput) error
}

type get_product_ratingUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetProductRatingUsecase creates a new instance
func NewGetProductRatingUsecase() GetProductRatingUsecase {
	return &get_product_ratingUsecaseImpl{}
}

// Execute executes GetProductRating
func (u *get_product_ratingUsecaseImpl) Execute(ctx context.Context, input GetProductRatingInput) error {
	// TODO: Implement business logic

	return nil
}
