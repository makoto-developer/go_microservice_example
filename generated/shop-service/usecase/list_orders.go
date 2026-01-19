package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type ListOrdersInput struct {
	ShopID uuid.UUID
	Limit  int
	Offset int
}

type ListOrdersOutput struct {
	Orders []*domain.Order
	Total  int
}

type ListOrdersUsecase interface {
	Execute(ctx context.Context, input ListOrdersInput) (*ListOrdersOutput, error)
}

type listOrdersUsecaseImpl struct {
	orderRepo domain.OrderRepository
}

func NewListOrdersUsecase(orderRepo domain.OrderRepository) ListOrdersUsecase {
	return &listOrdersUsecaseImpl{
		orderRepo: orderRepo,
	}
}

func (u *listOrdersUsecaseImpl) Execute(ctx context.Context, input ListOrdersInput) (*ListOrdersOutput, error) {
	limit := input.Limit
	if limit == 0 {
		limit = 20
	}

	orders, err := u.orderRepo.List(ctx, limit, input.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	return &ListOrdersOutput{
		Orders: orders,
		Total:  len(orders),
	}, nil
}
