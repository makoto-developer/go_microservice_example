package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// RejectShopInput represents input for RejectShop
type RejectShopInput struct {
	ShopId uuid.UUID
	AdminId uuid.UUID
	Reason string
}

// RejectShopUsecase defines the interface for RejectShop
type RejectShopUsecase interface {
	Execute(ctx context.Context, input RejectShopInput) error
}

type reject_shopUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewRejectShopUsecase creates a new instance
func NewRejectShopUsecase() RejectShopUsecase {
	return &reject_shopUsecaseImpl{}
}

// Execute executes RejectShop
func (u *reject_shopUsecaseImpl) Execute(ctx context.Context, input RejectShopInput) error {
	// TODO: Implement business logic

	return nil
}
