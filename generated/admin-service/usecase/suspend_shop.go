package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SuspendShopInput represents input for SuspendShop
type SuspendShopInput struct {
	ShopId uuid.UUID
	AdminId uuid.UUID
	Reason string
}

// SuspendShopUsecase defines the interface for SuspendShop
type SuspendShopUsecase interface {
	Execute(ctx context.Context, input SuspendShopInput) error
}

type suspend_shopUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewSuspendShopUsecase creates a new instance
func NewSuspendShopUsecase() SuspendShopUsecase {
	return &suspend_shopUsecaseImpl{}
}

// Execute executes SuspendShop
func (u *suspend_shopUsecaseImpl) Execute(ctx context.Context, input SuspendShopInput) error {
	// TODO: Implement business logic

	return nil
}
