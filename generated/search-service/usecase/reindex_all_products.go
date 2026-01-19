package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ReindexAllProductsInput represents input for ReindexAllProducts
type ReindexAllProductsInput struct {
	Output {
	TotalIndexed int
	Errors []string
}

// ReindexAllProductsUsecase defines the interface for ReindexAllProducts
type ReindexAllProductsUsecase interface {
	Execute(ctx context.Context, input ReindexAllProductsInput) error
}

type reindex_all_productsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewReindexAllProductsUsecase creates a new instance
func NewReindexAllProductsUsecase() ReindexAllProductsUsecase {
	return &reindex_all_productsUsecaseImpl{}
}

// Execute executes ReindexAllProducts
func (u *reindex_all_productsUsecaseImpl) Execute(ctx context.Context, input ReindexAllProductsInput) error {
	// TODO: Implement business logic

	return nil
}
