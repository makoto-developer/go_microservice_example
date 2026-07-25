package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type MergeGuestCartInput struct {
	CustomerID uuid.UUID
	SessionID  string
}

type MergeGuestCartOutput struct {
	MergedCount int
}

type MergeGuestCartUsecase interface {
	Execute(ctx context.Context, input MergeGuestCartInput) (MergeGuestCartOutput, error)
}

type mergeGuestCartUsecase struct {
	cartRepo repository.CartRepository
}

func NewMergeGuestCartUsecase(cartRepo repository.CartRepository) MergeGuestCartUsecase {
	return &mergeGuestCartUsecase{cartRepo: cartRepo}
}

func (u *mergeGuestCartUsecase) Execute(ctx context.Context, input MergeGuestCartInput) (MergeGuestCartOutput, error) {
	guestItems, err := u.cartRepo.GetBySessionID(ctx, input.SessionID)
	if err != nil {
		return MergeGuestCartOutput{}, err
	}

	mergedCount := 0
	for _, guestItem := range guestItems {
		cartItem := domain.NewCartItem(input.CustomerID, guestItem.ProductID, guestItem.VariationID, guestItem.Quantity)
		if err := u.cartRepo.AddItem(ctx, cartItem); err != nil {
			continue
		}
		mergedCount++
	}

	return MergeGuestCartOutput{MergedCount: mergedCount}, nil
}
