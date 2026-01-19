package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DeleteProductIndexInput represents input for DeleteProductIndex
type DeleteProductIndexInput struct {
	ProductId uuid.UUID
}

// DeleteProductIndexUsecase defines the interface for DeleteProductIndex
type DeleteProductIndexUsecase interface {
	Execute(ctx context.Context, input DeleteProductIndexInput) error
}

type delete_product_indexUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewDeleteProductIndexUsecase creates a new instance
func NewDeleteProductIndexUsecase() DeleteProductIndexUsecase {
	return &delete_product_indexUsecaseImpl{}
}

// Execute executes DeleteProductIndex
func (u *delete_product_indexUsecaseImpl) Execute(ctx context.Context, input DeleteProductIndexInput) error {
	// TODO: Implement business logic

	return nil
}
