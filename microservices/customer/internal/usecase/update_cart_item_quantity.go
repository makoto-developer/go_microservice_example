package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type UpdateCartItemQuantityInput struct {
	CustomerID uuid.UUID
	CartItemID uuid.UUID
	Quantity   int
}

type UpdateCartItemQuantityOutput struct {
	CartItem *domain.CartItem
}

type UpdateCartItemQuantityUsecase interface {
	Execute(ctx context.Context, input UpdateCartItemQuantityInput) (UpdateCartItemQuantityOutput, error)
}

type updateCartItemQuantityUsecase struct {
	cartRepo repository.CartRepository
}

func NewUpdateCartItemQuantityUsecase(cartRepo repository.CartRepository) UpdateCartItemQuantityUsecase {
	return &updateCartItemQuantityUsecase{cartRepo: cartRepo}
}

func (u *updateCartItemQuantityUsecase) Execute(ctx context.Context, input UpdateCartItemQuantityInput) (UpdateCartItemQuantityOutput, error) {
	if input.Quantity <= 0 {
		return UpdateCartItemQuantityOutput{}, domain.ErrInvalidQuantity
	}

	cartItem, err := u.cartRepo.GetByID(ctx, input.CartItemID)
	if err != nil {
		return UpdateCartItemQuantityOutput{}, err
	}

	if cartItem.CustomerID != input.CustomerID {
		return UpdateCartItemQuantityOutput{}, domain.ErrCartItemNotFound
	}

	if err := u.cartRepo.UpdateQuantity(ctx, input.CartItemID, input.Quantity); err != nil {
		return UpdateCartItemQuantityOutput{}, err
	}

	cartItem.Quantity = input.Quantity

	return UpdateCartItemQuantityOutput{CartItem: cartItem}, nil
}
