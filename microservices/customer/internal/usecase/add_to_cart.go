package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type AddToCartInput struct {
	CustomerID  uuid.UUID
	ProductID   uuid.UUID
	VariationID *uuid.UUID
	Quantity    int
}

type AddToCartOutput struct {
	CartItemID    uuid.UUID
	TotalQuantity int
}

type AddToCartUsecase interface {
	Execute(ctx context.Context, input AddToCartInput) (AddToCartOutput, error)
}

type addToCartUsecase struct {
	cartRepo repository.CartRepository
}

func NewAddToCartUsecase(cartRepo repository.CartRepository) AddToCartUsecase {
	return &addToCartUsecase{cartRepo: cartRepo}
}

func (u *addToCartUsecase) Execute(ctx context.Context, input AddToCartInput) (AddToCartOutput, error) {
	if input.Quantity <= 0 {
		return AddToCartOutput{}, domain.ErrInvalidQuantity
	}

	cartItem := domain.NewCartItem(input.CustomerID, input.ProductID, input.VariationID, input.Quantity)

	if err := u.cartRepo.AddItem(ctx, cartItem); err != nil {
		return AddToCartOutput{}, err
	}

	return AddToCartOutput{
		CartItemID:    cartItem.ID,
		TotalQuantity: input.Quantity,
	}, nil
}
