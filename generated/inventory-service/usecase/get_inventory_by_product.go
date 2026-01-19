package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetInventoryByProductInput represents input for GetInventoryByProduct
type GetInventoryByProductInput struct {
	ProductId uuid.UUID
	VariationId uuid.UUID
}

// GetInventoryByProductUsecase defines the interface for GetInventoryByProduct
type GetInventoryByProductUsecase interface {
	Execute(ctx context.Context, input GetInventoryByProductInput) error
}

type get_inventory_by_productUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetInventoryByProductUsecase creates a new instance
func NewGetInventoryByProductUsecase() GetInventoryByProductUsecase {
	return &get_inventory_by_productUsecaseImpl{}
}

// Execute executes GetInventoryByProduct
func (u *get_inventory_by_productUsecaseImpl) Execute(ctx context.Context, input GetInventoryByProductInput) error {
	// TODO: Implement business logic

	return nil
}
