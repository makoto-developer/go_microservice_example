package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository"
)

type GetCartInput struct {
	CustomerID uuid.UUID
}

type GetCartOutput struct {
	CartItems []*domain.CartItem
}

type GetCartUsecase interface {
	Execute(ctx context.Context, input GetCartInput) (GetCartOutput, error)
}

type getCartUsecase struct {
	cartRepo repository.CartRepository
}

func NewGetCartUsecase(cartRepo repository.CartRepository) GetCartUsecase {
	return &getCartUsecase{cartRepo: cartRepo}
}

func (u *getCartUsecase) Execute(ctx context.Context, input GetCartInput) (GetCartOutput, error) {
	items, err := u.cartRepo.GetByCustomerID(ctx, input.CustomerID)
	if err != nil {
		return GetCartOutput{}, err
	}

	return GetCartOutput{CartItems: items}, nil
}
