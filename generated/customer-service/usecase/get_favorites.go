package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetFavoritesInput represents input for GetFavorites
type GetFavoritesInput struct {
	CustomerId uuid.UUID
	SortBy string
	SortOrder string
}

// GetFavoritesUsecase defines the interface for GetFavorites
type GetFavoritesUsecase interface {
	Execute(ctx context.Context, input GetFavoritesInput) error
}

type get_favoritesUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetFavoritesUsecase creates a new instance
func NewGetFavoritesUsecase() GetFavoritesUsecase {
	return &get_favoritesUsecaseImpl{}
}

// Execute executes GetFavorites
func (u *get_favoritesUsecaseImpl) Execute(ctx context.Context, input GetFavoritesInput) error {
	// TODO: Implement business logic

	return nil
}
