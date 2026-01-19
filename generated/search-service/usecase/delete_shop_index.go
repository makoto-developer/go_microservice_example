package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DeleteShopIndexInput represents input for DeleteShopIndex
type DeleteShopIndexInput struct {
	ShopId uuid.UUID
}

// DeleteShopIndexUsecase defines the interface for DeleteShopIndex
type DeleteShopIndexUsecase interface {
	Execute(ctx context.Context, input DeleteShopIndexInput) error
}

type delete_shop_indexUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewDeleteShopIndexUsecase creates a new instance
func NewDeleteShopIndexUsecase() DeleteShopIndexUsecase {
	return &delete_shop_indexUsecaseImpl{}
}

// Execute executes DeleteShopIndex
func (u *delete_shop_indexUsecaseImpl) Execute(ctx context.Context, input DeleteShopIndexInput) error {
	// TODO: Implement business logic

	return nil
}
