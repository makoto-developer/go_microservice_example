package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type RemoveFromCartInput struct {
	CustomerID uuid.UUID
	CartItemID uuid.UUID
}

type RemoveFromCartOutput struct {
	Success bool
}

type RemoveFromCartUsecase interface {
	Execute(ctx context.Context, input RemoveFromCartInput) (RemoveFromCartOutput, error)
}

type removeFromCartUsecase struct {
	cartRepo repository.CartRepository
}

func NewRemoveFromCartUsecase(cartRepo repository.CartRepository) RemoveFromCartUsecase {
	return &removeFromCartUsecase{cartRepo: cartRepo}
}

func (u *removeFromCartUsecase) Execute(ctx context.Context, input RemoveFromCartInput) (RemoveFromCartOutput, error) {
	cartItem, err := u.cartRepo.GetByID(ctx, input.CartItemID)
	if err != nil {
		return RemoveFromCartOutput{}, err
	}

	if cartItem.CustomerID != input.CustomerID {
		return RemoveFromCartOutput{}, domain.ErrCartItemNotFound
	}

	if err := u.cartRepo.RemoveItem(ctx, input.CartItemID); err != nil {
		return RemoveFromCartOutput{}, err
	}

	return RemoveFromCartOutput{Success: true}, nil
}
