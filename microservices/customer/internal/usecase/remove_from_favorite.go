package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type RemoveFromFavoriteInput struct {
	CustomerID uuid.UUID
	ProductID  uuid.UUID
}

type RemoveFromFavoriteOutput struct {
	Success bool
}

type RemoveFromFavoriteUsecase interface {
	Execute(ctx context.Context, input RemoveFromFavoriteInput) (RemoveFromFavoriteOutput, error)
}

type removeFromFavoriteUsecase struct {
	favoriteRepo repository.FavoriteRepository
}

func NewRemoveFromFavoriteUsecase(favoriteRepo repository.FavoriteRepository) RemoveFromFavoriteUsecase {
	return &removeFromFavoriteUsecase{favoriteRepo: favoriteRepo}
}

func (u *removeFromFavoriteUsecase) Execute(ctx context.Context, input RemoveFromFavoriteInput) (RemoveFromFavoriteOutput, error) {
	if err := u.favoriteRepo.Remove(ctx, input.CustomerID, input.ProductID); err != nil {
		return RemoveFromFavoriteOutput{}, err
	}

	return RemoveFromFavoriteOutput{Success: true}, nil
}
