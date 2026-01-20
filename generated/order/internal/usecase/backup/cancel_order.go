package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/repository"
)

type CancelOrderInput struct {
	OrderID uuid.UUID
}

type CancelOrderOutput struct {
	Cancelled bool
}

type CancelOrderUsecase interface {
	Execute(ctx context.Context, input CancelOrderInput) (CancelOrderOutput, error)
}

type cancelOrderUsecaseImpl struct {
	orderRepo repository.OrderRepository
}

func NewCancelOrderUsecase(orderRepo repository.OrderRepository) CancelOrderUsecase {
	return &cancelOrderUsecaseImpl{
		orderRepo: orderRepo,
	}
}

func (u *cancelOrderUsecaseImpl) Execute(ctx context.Context, input CancelOrderInput) (CancelOrderOutput, error) {
	// Get order
	order, err := u.orderRepo.GetByID(ctx, input.OrderID)
	if err != nil {
		return CancelOrderOutput{}, err
	}

	// Cancel order
	err = order.Cancel()
	if err != nil {
		return CancelOrderOutput{}, err
	}

	// Update repository
	err = u.orderRepo.Cancel(ctx, order.ID)
	if err != nil {
		return CancelOrderOutput{}, err
	}

	return CancelOrderOutput{Cancelled: true}, nil
}
