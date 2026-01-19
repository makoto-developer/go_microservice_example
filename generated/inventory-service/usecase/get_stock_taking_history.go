package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetStockTakingHistoryInput represents input for GetStockTakingHistory
type GetStockTakingHistoryInput struct {
	ShopId uuid.UUID
	DateFrom date
	DateTo date
}

// GetStockTakingHistoryUsecase defines the interface for GetStockTakingHistory
type GetStockTakingHistoryUsecase interface {
	Execute(ctx context.Context, input GetStockTakingHistoryInput) error
}

type get_stock_taking_historyUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetStockTakingHistoryUsecase creates a new instance
func NewGetStockTakingHistoryUsecase() GetStockTakingHistoryUsecase {
	return &get_stock_taking_historyUsecaseImpl{}
}

// Execute executes GetStockTakingHistory
func (u *get_stock_taking_historyUsecaseImpl) Execute(ctx context.Context, input GetStockTakingHistoryInput) error {
	// TODO: Implement business logic

	return nil
}
