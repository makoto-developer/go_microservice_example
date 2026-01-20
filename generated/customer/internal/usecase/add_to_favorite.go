package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository"
)

type AddToFavoriteInput struct {
	CustomerID      uuid.UUID
	ProductID       uuid.UUID
	NotifyOnRestock bool
}

type AddToFavoriteOutput struct {
	FavoriteID uuid.UUID
}

type AddToFavoriteUsecase interface {
	Execute(ctx context.Context, input AddToFavoriteInput) (AddToFavoriteOutput, error)
}

type addToFavoriteUsecase struct {
	favoriteRepo repository.FavoriteRepository
}

func NewAddToFavoriteUsecase(favoriteRepo repository.FavoriteRepository) AddToFavoriteUsecase {
	return &addToFavoriteUsecase{favoriteRepo: favoriteRepo}
}

func (u *addToFavoriteUsecase) Execute(ctx context.Context, input AddToFavoriteInput) (AddToFavoriteOutput, error) {
	exists, err := u.favoriteRepo.Exists(ctx, input.CustomerID, input.ProductID)
	if err != nil {
		return AddToFavoriteOutput{}, err
	}
	if exists {
		return AddToFavoriteOutput{}, domain.ErrAlreadyFavorited
	}

	favorite := domain.NewFavorite(input.CustomerID, input.ProductID)

	if err := u.favoriteRepo.Add(ctx, favorite); err != nil {
		return AddToFavoriteOutput{}, err
	}

	return AddToFavoriteOutput{FavoriteID: favorite.ID}, nil
}
