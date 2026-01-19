package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
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
	Message string
}

type UpdateOrderStatusUsecase interface {
	Execute(ctx context.Context, input UpdateOrderStatusInput) (*UpdateOrderStatusOutput, error)
}

type updateOrderStatusUsecaseImpl struct {
	orderRepo domain.OrderRepository
}

func NewUpdateOrderStatusUsecase(orderRepo domain.OrderRepository) UpdateOrderStatusUsecase {
	return &updateOrderStatusUsecaseImpl{
		orderRepo: orderRepo,
	}
}

func (u *updateOrderStatusUsecaseImpl) Execute(ctx context.Context, input UpdateOrderStatusInput) (*UpdateOrderStatusOutput, error) {
	order, err := u.orderRepo.FindByID(ctx, input.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	if order.ShopId != input.ShopID {
		return nil, fmt.Errorf("order does not belong to shop")
	}

	order.Status = input.NewStatus
	if input.TrackingNumber != nil {
		order.TrackingNumber = input.TrackingNumber
	}
	if input.Carrier != nil {
		order.Carrier = input.Carrier
	}
	order.UpdatedAt = time.Now()

	if err := u.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return &UpdateOrderStatusOutput{
		OrderID: order.Id,
		Status:  order.Status,
		Message: "Order status updated successfully",
	}, nil
}
