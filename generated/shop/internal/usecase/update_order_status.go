package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type UpdateOrderStatusInput struct {
	OrderID        uuid.UUID
	ShopID         uuid.UUID
	NewStatus      domain.OrderStatus
	TrackingNumber *string
	Carrier        *domain.Carrier
}

type UpdateOrderStatusOutput struct {
	OrderID uuid.UUID
	Status  domain.OrderStatus
}

type UpdateOrderStatusUsecase interface {
	Execute(ctx context.Context, input UpdateOrderStatusInput) (UpdateOrderStatusOutput, error)
}

type updateOrderStatusUsecase struct {
	orderRepo repository.OrderRepository
}

func NewUpdateOrderStatusUsecase(orderRepo repository.OrderRepository) UpdateOrderStatusUsecase {
	return &updateOrderStatusUsecase{orderRepo: orderRepo}
}

func (u *updateOrderStatusUsecase) Execute(ctx context.Context, input UpdateOrderStatusInput) (UpdateOrderStatusOutput, error) {
	order, err := u.orderRepo.GetByID(ctx, input.OrderID)
	if err != nil {
		return UpdateOrderStatusOutput{}, err
	}

	if order.ShopID != input.ShopID {
		return UpdateOrderStatusOutput{}, domain.ErrUnauthorizedAccess
	}

	if !order.CanUpdateStatus(input.NewStatus) {
		return UpdateOrderStatusOutput{}, domain.ErrInvalidStatusTransition
	}

	order.Status = input.NewStatus
	order.TrackingNumber = input.TrackingNumber
	order.Carrier = input.Carrier
	order.UpdatedAt = time.Now()

	if err := u.orderRepo.Update(ctx, order); err != nil {
		return UpdateOrderStatusOutput{}, err
	}

	return UpdateOrderStatusOutput{
		OrderID: order.ID,
		Status:  order.Status,
	}, nil
}
