package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreateReorderInput represents input for CreateReorder
type CreateReorderInput struct {
	CustomerId uuid.UUID
	OriginalOrderId uuid.UUID
}

// CreateReorderUsecase defines the interface for CreateReorder
type CreateReorderUsecase interface {
	Execute(ctx context.Context, input CreateReorderInput) error
}

type create_reorderUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewCreateReorderUsecase creates a new instance
func NewCreateReorderUsecase() CreateReorderUsecase {
	return &create_reorderUsecaseImpl{}
}

// Execute executes CreateReorder
func (u *create_reorderUsecaseImpl) Execute(ctx context.Context, input CreateReorderInput) error {
	// TODO: Implement business logic

	return nil
}
