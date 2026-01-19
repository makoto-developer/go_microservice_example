package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// RemoveFromFavoriteInput represents input for RemoveFromFavorite
type RemoveFromFavoriteInput struct {
	FavoriteId uuid.UUID
	CustomerId uuid.UUID
}

// RemoveFromFavoriteUsecase defines the interface for RemoveFromFavorite
type RemoveFromFavoriteUsecase interface {
	Execute(ctx context.Context, input RemoveFromFavoriteInput) error
}

type remove_from_favoriteUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewRemoveFromFavoriteUsecase creates a new instance
func NewRemoveFromFavoriteUsecase() RemoveFromFavoriteUsecase {
	return &remove_from_favoriteUsecaseImpl{}
}

// Execute executes RemoveFromFavorite
func (u *remove_from_favoriteUsecaseImpl) Execute(ctx context.Context, input RemoveFromFavoriteInput) error {
	// TODO: Implement business logic

	return nil
}
