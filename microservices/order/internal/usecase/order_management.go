package usecase

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/repository"
)

type CreateOrderInput struct {
	CustomerID  uuid.UUID
	AddressID   uuid.UUID
	Items       []OrderItemInput
	ShippingFee int64
}

type OrderItemInput struct {
	ProductID   uuid.UUID
	VariationID *uuid.UUID
	Quantity    int
	UnitPrice   int64
}

type OrderManagementUsecase interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (uuid.UUID, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (*domain.Order, error)
	GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderItem, error)
	CancelOrder(ctx context.Context, orderID uuid.UUID) error
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error
}

type orderManagementUsecase struct {
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
}

func NewOrderManagementUsecase(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
) OrderManagementUsecase {
	return &orderManagementUsecase{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

func (u *orderManagementUsecase) CreateOrder(ctx context.Context, input CreateOrderInput) (uuid.UUID, error) {
	var totalAmount int64
	for _, item := range input.Items {
		totalAmount += item.UnitPrice * int64(item.Quantity)
	}
	totalAmount += input.ShippingFee
	
	order := &domain.Order{
		ID:          uuid.New(),
		CustomerID:  input.CustomerID,
		OrderNumber: fmt.Sprintf("ORD-%d", time.Now().Unix()),
		Status:      domain.OrderStatusPending,
		TotalAmount: totalAmount,
		ShippingFee: input.ShippingFee,
		AddressID:   input.AddressID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	if err := u.orderRepo.Create(ctx, order); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	for _, itemInput := range input.Items {
		item := &domain.OrderItem{
			ID:          uuid.New(),
			OrderID:     order.ID,
			ProductID:   itemInput.ProductID,
			VariationID: itemInput.VariationID,
			Quantity:    itemInput.Quantity,
			UnitPrice:   itemInput.UnitPrice,
			Subtotal:    itemInput.UnitPrice * int64(itemInput.Quantity),
			CreatedAt:   time.Now(),
		}
		
		if err := u.orderItemRepo.Create(ctx, item); err != nil {
			return uuid.Nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}
	
	return order.ID, nil
}

func (u *orderManagementUsecase) GetOrder(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

func (u *orderManagementUsecase) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderItem, error) {
	items, err := u.orderItemRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	return items, nil
}

func (u *orderManagementUsecase) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	if err := u.orderRepo.UpdateStatus(ctx, orderID, domain.OrderStatusCancelled); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	return nil
}

func (u *orderManagementUsecase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	if err := u.orderRepo.UpdateStatus(ctx, orderID, status); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}
