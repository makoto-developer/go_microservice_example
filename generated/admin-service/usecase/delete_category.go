package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DeleteCategoryInput represents input for DeleteCategory
type DeleteCategoryInput struct {
	AdminId uuid.UUID
	CategoryId uuid.UUID
}

// DeleteCategoryUsecase defines the interface for DeleteCategory
type DeleteCategoryUsecase interface {
	Execute(ctx context.Context, input DeleteCategoryInput) error
}

type delete_categoryUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewDeleteCategoryUsecase creates a new instance
func NewDeleteCategoryUsecase() DeleteCategoryUsecase {
	return &delete_categoryUsecaseImpl{}
}

// Execute executes DeleteCategory
func (u *delete_categoryUsecaseImpl) Execute(ctx context.Context, input DeleteCategoryInput) error {
	// TODO: Implement business logic

	return nil
}
