package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ExportOrdersToCSVInput represents input for ExportOrdersToCSV
type ExportOrdersToCSVInput struct {
	ShopId uuid.UUID
	DateFrom date
	DateTo date
	StatusFilter []OrderStatus
}

// ExportOrdersToCSVUsecase defines the interface for ExportOrdersToCSV
type ExportOrdersToCSVUsecase interface {
	Execute(ctx context.Context, input ExportOrdersToCSVInput) error
}

type export_orders_to_c_s_vUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewExportOrdersToCSVUsecase creates a new instance
func NewExportOrdersToCSVUsecase() ExportOrdersToCSVUsecase {
	return &export_orders_to_c_s_vUsecaseImpl{}
}

// Execute executes ExportOrdersToCSV
func (u *export_orders_to_c_s_vUsecaseImpl) Execute(ctx context.Context, input ExportOrdersToCSVInput) error {
	// TODO: Implement business logic

	return nil
}
