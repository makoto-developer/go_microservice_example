package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetAllShopsInput represents input for GetAllShops
type GetAllShopsInput struct {
	ApprovalStatus string
	StatusFilter string
	Category string
	DateFrom date
	DateTo date
	SortBy string
	Page int
	PageSize int
}

// GetAllShopsUsecase defines the interface for GetAllShops
type GetAllShopsUsecase interface {
	Execute(ctx context.Context, input GetAllShopsInput) error
}

type get_all_shopsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetAllShopsUsecase creates a new instance
func NewGetAllShopsUsecase() GetAllShopsUsecase {
	return &get_all_shopsUsecaseImpl{}
}

// Execute executes GetAllShops
func (u *get_all_shopsUsecaseImpl) Execute(ctx context.Context, input GetAllShopsInput) error {
	// TODO: Implement business logic

	return nil
}
