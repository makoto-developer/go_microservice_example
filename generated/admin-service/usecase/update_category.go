package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpdateCategoryInput represents input for UpdateCategory
type UpdateCategoryInput struct {
	AdminId uuid.UUID
	CategoryId uuid.UUID
	Name string
	DisplayOrder int
}

// UpdateCategoryUsecase defines the interface for UpdateCategory
type UpdateCategoryUsecase interface {
	Execute(ctx context.Context, input UpdateCategoryInput) error
}

type update_categoryUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUpdateCategoryUsecase creates a new instance
func NewUpdateCategoryUsecase() UpdateCategoryUsecase {
	return &update_categoryUsecaseImpl{}
}

// Execute executes UpdateCategory
func (u *update_categoryUsecaseImpl) Execute(ctx context.Context, input UpdateCategoryInput) error {
	// TODO: Implement business logic

	return nil
}
