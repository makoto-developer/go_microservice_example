package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreateCategoryInput represents input for CreateCategory
type CreateCategoryInput struct {
	AdminId uuid.UUID
	Name string
	ParentId uuid.UUID
	DisplayOrder int
}

// CreateCategoryUsecase defines the interface for CreateCategory
type CreateCategoryUsecase interface {
	Execute(ctx context.Context, input CreateCategoryInput) error
}

type create_categoryUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewCreateCategoryUsecase creates a new instance
func NewCreateCategoryUsecase() CreateCategoryUsecase {
	return &create_categoryUsecaseImpl{}
}

// Execute executes CreateCategory
func (u *create_categoryUsecaseImpl) Execute(ctx context.Context, input CreateCategoryInput) error {
	// TODO: Implement business logic

	return nil
}
