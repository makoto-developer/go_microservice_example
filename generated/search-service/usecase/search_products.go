package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SearchProductsInput represents input for SearchProducts
type SearchProductsInput struct {
	Keyword string
	Category string
	MinPrice decimal.Decimal
	MaxPrice decimal.Decimal
	MinRating decimal.Decimal
	StockStatus StockStatus
	ShopId uuid.UUID
	SortBy SortBy
	Page int
	PageSize int
}

// SearchProductsUsecase defines the interface for SearchProducts
type SearchProductsUsecase interface {
	Execute(ctx context.Context, input SearchProductsInput) error
}

type search_productsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewSearchProductsUsecase creates a new instance
func NewSearchProductsUsecase() SearchProductsUsecase {
	return &search_productsUsecaseImpl{}
}

// Execute executes SearchProducts
func (u *search_productsUsecaseImpl) Execute(ctx context.Context, input SearchProductsInput) error {
	// TODO: Implement business logic

	return nil
}
