package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ActivateShopInput represents input for ActivateShop
type ActivateShopInput struct {
	ShopId uuid.UUID
	AdminId uuid.UUID
	Reason string
}

// ActivateShopUsecase defines the interface for ActivateShop
type ActivateShopUsecase interface {
	Execute(ctx context.Context, input ActivateShopInput) error
}

type activate_shopUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewActivateShopUsecase creates a new instance
func NewActivateShopUsecase() ActivateShopUsecase {
	return &activate_shopUsecaseImpl{}
}

// Execute executes ActivateShop
func (u *activate_shopUsecaseImpl) Execute(ctx context.Context, input ActivateShopInput) error {
	// TODO: Implement business logic

	return nil
}
