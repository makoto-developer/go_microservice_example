package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type GetOrderDetailInput struct {
	OrderID uuid.UUID
	ShopID  uuid.UUID
}

type GetOrderDetailOutput struct {
	Order      *domain.Order
	OrderItems []*domain.OrderItem
}

type GetOrderDetailUsecase interface {
	Execute(ctx context.Context, input GetOrderDetailInput) (*GetOrderDetailOutput, error)
}

type getOrderDetailUsecaseImpl struct {
	orderRepo     domain.OrderRepository
	orderItemRepo domain.OrderItemRepository
}

func NewGetOrderDetailUsecase(
	orderRepo domain.OrderRepository,
	orderItemRepo domain.OrderItemRepository,
) GetOrderDetailUsecase {
	return &getOrderDetailUsecaseImpl{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

func (u *getOrderDetailUsecaseImpl) Execute(ctx context.Context, input GetOrderDetailInput) (*GetOrderDetailOutput, error) {
	order, err := u.orderRepo.FindByID(ctx, input.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	if order.ShopId != input.ShopID {
		return nil, fmt.Errorf("order does not belong to shop")
	}

	orderItems, err := u.orderItemRepo.List(ctx, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to list order items: %w", err)
	}

	return &GetOrderDetailOutput{
		Order:      order,
		OrderItems: orderItems,
	}, nil
}
