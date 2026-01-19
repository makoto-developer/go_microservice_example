package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetInventoryHistoryInput represents input for GetInventoryHistory
type GetInventoryHistoryInput struct {
	InventoryId uuid.UUID
	DateFrom date
	DateTo date
	ChangeTypes []ChangeType
	Page int
	PageSize int
}

// GetInventoryHistoryUsecase defines the interface for GetInventoryHistory
type GetInventoryHistoryUsecase interface {
	Execute(ctx context.Context, input GetInventoryHistoryInput) error
}

type get_inventory_historyUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetInventoryHistoryUsecase creates a new instance
func NewGetInventoryHistoryUsecase() GetInventoryHistoryUsecase {
	return &get_inventory_historyUsecaseImpl{}
}

// Execute executes GetInventoryHistory
func (u *get_inventory_historyUsecaseImpl) Execute(ctx context.Context, input GetInventoryHistoryInput) error {
	// TODO: Implement business logic

	return nil
}
