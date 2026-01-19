package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetCategoriesInput represents input for GetCategories
type GetCategoriesInput struct {
	Output {
	Categories []Category
}

// GetCategoriesUsecase defines the interface for GetCategories
type GetCategoriesUsecase interface {
	Execute(ctx context.Context, input GetCategoriesInput) error
}

type get_categoriesUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetCategoriesUsecase creates a new instance
func NewGetCategoriesUsecase() GetCategoriesUsecase {
	return &get_categoriesUsecaseImpl{}
}

// Execute executes GetCategories
func (u *get_categoriesUsecaseImpl) Execute(ctx context.Context, input GetCategoriesInput) error {
	// TODO: Implement business logic

	return nil
}
