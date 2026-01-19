package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// AddToFavoriteInput represents input for AddToFavorite
type AddToFavoriteInput struct {
	CustomerId uuid.UUID
	ProductId uuid.UUID
	NotifyOnRestock bool
}

// AddToFavoriteUsecase defines the interface for AddToFavorite
type AddToFavoriteUsecase interface {
	Execute(ctx context.Context, input AddToFavoriteInput) error
}

type add_to_favoriteUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewAddToFavoriteUsecase creates a new instance
func NewAddToFavoriteUsecase() AddToFavoriteUsecase {
	return &add_to_favoriteUsecaseImpl{}
}

// Execute executes AddToFavorite
func (u *add_to_favoriteUsecaseImpl) Execute(ctx context.Context, input AddToFavoriteInput) error {
	// TODO: Implement business logic

	return nil
}
