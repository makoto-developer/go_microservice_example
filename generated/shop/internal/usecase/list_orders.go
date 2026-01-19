package usecase

import (
	"context"

	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type ListOrdersInput struct {
	Filter repository.OrderFilter
}

type ListOrdersOutput struct {
	Orders     []*domain.Order
	TotalCount int
}

type ListOrdersUsecase interface {
	Execute(ctx context.Context, input ListOrdersInput) (ListOrdersOutput, error)
}

type listOrdersUsecase struct {
	orderRepo repository.OrderRepository
}

func NewListOrdersUsecase(orderRepo repository.OrderRepository) ListOrdersUsecase {
	return &listOrdersUsecase{orderRepo: orderRepo}
}

func (u *listOrdersUsecase) Execute(ctx context.Context, input ListOrdersInput) (ListOrdersOutput, error) {
	orders, err := u.orderRepo.List(ctx, input.Filter)
	if err != nil {
		return ListOrdersOutput{}, err
	}

	return ListOrdersOutput{
		Orders:     orders,
		TotalCount: len(orders),
	}, nil
}
