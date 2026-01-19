package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// IndexShopInput represents input for IndexShop
type IndexShopInput struct {
	ShopId uuid.UUID
	ShopName string
	Description string
	Categories []string
	AverageRating decimal.Decimal
	ProductCount int
	LogoUrl string
}

// IndexShopUsecase defines the interface for IndexShop
type IndexShopUsecase interface {
	Execute(ctx context.Context, input IndexShopInput) error
}

type index_shopUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewIndexShopUsecase creates a new instance
func NewIndexShopUsecase() IndexShopUsecase {
	return &index_shopUsecaseImpl{}
}

// Execute executes IndexShop
func (u *index_shopUsecaseImpl) Execute(ctx context.Context, input IndexShopInput) error {
	// TODO: Implement business logic

	return nil
}
