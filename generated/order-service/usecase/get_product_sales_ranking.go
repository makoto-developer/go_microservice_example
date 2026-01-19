package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetProductSalesRankingInput represents input for GetProductSalesRanking
type GetProductSalesRankingInput struct {
	ShopId uuid.UUID
	DateFrom date
	DateTo date
	Limit int
}

// GetProductSalesRankingUsecase defines the interface for GetProductSalesRanking
type GetProductSalesRankingUsecase interface {
	Execute(ctx context.Context, input GetProductSalesRankingInput) error
}

type get_product_sales_rankingUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetProductSalesRankingUsecase creates a new instance
func NewGetProductSalesRankingUsecase() GetProductSalesRankingUsecase {
	return &get_product_sales_rankingUsecaseImpl{}
}

// Execute executes GetProductSalesRanking
func (u *get_product_sales_rankingUsecaseImpl) Execute(ctx context.Context, input GetProductSalesRankingInput) error {
	// TODO: Implement business logic

	return nil
}
