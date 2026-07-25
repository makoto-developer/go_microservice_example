package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type ListFavoritesInput struct {
	CustomerID uuid.UUID
}

type ListFavoritesOutput struct {
	Favorites []*domain.Favorite
}

type ListFavoritesUsecase interface {
	Execute(ctx context.Context, input ListFavoritesInput) (ListFavoritesOutput, error)
}

type listFavoritesUsecase struct {
	favoriteRepo repository.FavoriteRepository
}

func NewListFavoritesUsecase(favoriteRepo repository.FavoriteRepository) ListFavoritesUsecase {
	return &listFavoritesUsecase{favoriteRepo: favoriteRepo}
}

func (u *listFavoritesUsecase) Execute(ctx context.Context, input ListFavoritesInput) (ListFavoritesOutput, error) {
	favorites, err := u.favoriteRepo.List(ctx, input.CustomerID)
	if err != nil {
		return ListFavoritesOutput{}, err
	}

	return ListFavoritesOutput{Favorites: favorites}, nil
}
