package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// IndexProductInput represents input for IndexProduct
type IndexProductInput struct {
	ProductId uuid.UUID
	ProductName string
	Description string
	Category string
	Tags []string
	ShopId uuid.UUID
	ShopName string
	Price decimal.Decimal
	AverageRating decimal.Decimal
	ReviewCount int
	StockStatus StockStatus
	ImageUrl string
}

// IndexProductUsecase defines the interface for IndexProduct
type IndexProductUsecase interface {
	Execute(ctx context.Context, input IndexProductInput) error
}

type index_productUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewIndexProductUsecase creates a new instance
func NewIndexProductUsecase() IndexProductUsecase {
	return &index_productUsecaseImpl{}
}

// Execute executes IndexProduct
func (u *index_productUsecaseImpl) Execute(ctx context.Context, input IndexProductInput) error {
	// TODO: Implement business logic

	return nil
}
