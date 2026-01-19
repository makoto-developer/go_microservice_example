package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetSalesChartInput represents input for GetSalesChart
type GetSalesChartInput struct {
	DateFrom date
	DateTo date
	GroupBy string
}

// GetSalesChartUsecase defines the interface for GetSalesChart
type GetSalesChartUsecase interface {
	Execute(ctx context.Context, input GetSalesChartInput) error
}

type get_sales_chartUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetSalesChartUsecase creates a new instance
func NewGetSalesChartUsecase() GetSalesChartUsecase {
	return &get_sales_chartUsecaseImpl{}
}

// Execute executes GetSalesChart
func (u *get_sales_chartUsecaseImpl) Execute(ctx context.Context, input GetSalesChartInput) error {
	// TODO: Implement business logic

	return nil
}
