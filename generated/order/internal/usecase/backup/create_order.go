package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/repository"
)

type OrderItemInput struct {
	ProductID uuid.UUID
	ShopID    uuid.UUID
	Quantity  int
	Price     int
}

type CreateOrderInput struct {
	CustomerID uuid.UUID
	Items      []OrderItemInput
}

type CreateOrderOutput struct {
	OrderID     uuid.UUID
	OrderNumber string
	TotalAmount int
}

type CreateOrderUsecase interface {
	Execute(ctx context.Context, input CreateOrderInput) (CreateOrderOutput, error)
}

type createOrderUsecaseImpl struct {
	orderRepo repository.OrderRepository
}

func NewCreateOrderUsecase(orderRepo repository.OrderRepository) CreateOrderUsecase {
	return &createOrderUsecaseImpl{
		orderRepo: orderRepo,
	}
}

func (u *createOrderUsecaseImpl) Execute(ctx context.Context, input CreateOrderInput) (CreateOrderOutput, error) {
	// Calculate total amount
	totalAmount := 0
	for _, item := range input.Items {
		totalAmount += item.Price * item.Quantity
	}

	// Generate order number
	orderNumber := generateOrderNumber()

	// Create order
	order := domain.NewOrder(input.CustomerID, orderNumber, totalAmount)

	// Create order items
	items := make([]*domain.OrderItem, len(input.Items))
	for i, itemInput := range input.Items {
		items[i] = domain.NewOrderItem(
			order.ID,
			itemInput.ProductID,
			itemInput.ShopID,
			itemInput.Quantity,
			itemInput.Price,
		)
	}

	// Save to repository
	err := u.orderRepo.Create(ctx, order, items)
	if err != nil {
		return CreateOrderOutput{}, err
	}

	return CreateOrderOutput{
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		TotalAmount: order.TotalAmount,
	}, nil
}

func generateOrderNumber() string {
	// Format: ORD-20260118-XXXXXX
	now := time.Now()
	random := uuid.New().String()[:6]
	return fmt.Sprintf("ORD-%s-%s", now.Format("20060102"), random)
}
