package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpdateInventoryQuantityInput represents input for UpdateInventoryQuantity
type UpdateInventoryQuantityInput struct {
	InventoryId uuid.UUID
	ShopId uuid.UUID
	ChangeQuantity int
	ChangeType ChangeType
	Reason string
	Operator string
}

// UpdateInventoryQuantityUsecase defines the interface for UpdateInventoryQuantity
type UpdateInventoryQuantityUsecase interface {
	Execute(ctx context.Context, input UpdateInventoryQuantityInput) error
}

type update_inventory_quantityUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUpdateInventoryQuantityUsecase creates a new instance
func NewUpdateInventoryQuantityUsecase() UpdateInventoryQuantityUsecase {
	return &update_inventory_quantityUsecaseImpl{}
}

// Execute executes UpdateInventoryQuantity
func (u *update_inventory_quantityUsecaseImpl) Execute(ctx context.Context, input UpdateInventoryQuantityInput) error {
	// TODO: Implement business logic

	return nil
}
