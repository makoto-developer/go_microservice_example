package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SearchShopsInput represents input for SearchShops
type SearchShopsInput struct {
	Keyword string
	Category string
	MinRating decimal.Decimal
	SortBy string
	Page int
	PageSize int
}

// SearchShopsUsecase defines the interface for SearchShops
type SearchShopsUsecase interface {
	Execute(ctx context.Context, input SearchShopsInput) error
}

type search_shopsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewSearchShopsUsecase creates a new instance
func NewSearchShopsUsecase() SearchShopsUsecase {
	return &search_shopsUsecaseImpl{}
}

// Execute executes SearchShops
func (u *search_shopsUsecaseImpl) Execute(ctx context.Context, input SearchShopsInput) error {
	// TODO: Implement business logic

	return nil
}
