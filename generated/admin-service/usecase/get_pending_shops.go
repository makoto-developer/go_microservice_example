package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetPendingShopsInput represents input for GetPendingShops
type GetPendingShopsInput struct {
	Page int
	PageSize int
}

// GetPendingShopsUsecase defines the interface for GetPendingShops
type GetPendingShopsUsecase interface {
	Execute(ctx context.Context, input GetPendingShopsInput) error
}

type get_pending_shopsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetPendingShopsUsecase creates a new instance
func NewGetPendingShopsUsecase() GetPendingShopsUsecase {
	return &get_pending_shopsUsecaseImpl{}
}

// Execute executes GetPendingShops
func (u *get_pending_shopsUsecaseImpl) Execute(ctx context.Context, input GetPendingShopsInput) error {
	// TODO: Implement business logic

	return nil
}
