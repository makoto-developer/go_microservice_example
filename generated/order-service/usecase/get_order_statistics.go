package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetOrderStatisticsInput represents input for GetOrderStatistics
type GetOrderStatisticsInput struct {
	ShopId uuid.UUID
	DateFrom date
	DateTo date
	GroupBy string
}

// GetOrderStatisticsUsecase defines the interface for GetOrderStatistics
type GetOrderStatisticsUsecase interface {
	Execute(ctx context.Context, input GetOrderStatisticsInput) error
}

type get_order_statisticsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetOrderStatisticsUsecase creates a new instance
func NewGetOrderStatisticsUsecase() GetOrderStatisticsUsecase {
	return &get_order_statisticsUsecaseImpl{}
}

// Execute executes GetOrderStatistics
func (u *get_order_statisticsUsecaseImpl) Execute(ctx context.Context, input GetOrderStatisticsInput) error {
	// TODO: Implement business logic

	return nil
}
