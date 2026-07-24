package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/client"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/repository"
	"time"
)

// PaymentMethod は注文時に選択される支払い方法。
type PaymentMethod string

const (
	PaymentMethodCreditCard     PaymentMethod = "credit_card"
	PaymentMethodCashOnDelivery PaymentMethod = "cash_on_delivery"
)

type CreateOrderInput struct {
	CustomerID      uuid.UUID
	AddressID       uuid.UUID
	Items           []OrderItemInput
	ShippingFee     int64
	PaymentMethod   PaymentMethod // 未指定はクレジットカード扱い
	PaymentMethodID string
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
	paymentClient client.PaymentClient
}

func NewOrderManagementUsecase(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	paymentClient client.PaymentClient,
) OrderManagementUsecase {
	return &orderManagementUsecase{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		paymentClient: paymentClient,
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

	// 決済サービスで支払いを実行する。失敗した場合は注文をキャンセルする。
	// paymentClient が未設定(nil)の場合は決済をスキップする(ローカル起動の後方互換)。
	if u.paymentClient != nil {
		nextStatus, err := u.executePayment(ctx, order.ID, input, totalAmount)
		if err != nil {
			if cancelErr := u.orderRepo.UpdateStatus(ctx, order.ID, domain.OrderStatusCancelled); cancelErr != nil {
				return uuid.Nil, fmt.Errorf("payment failed (%v) and order cancellation also failed: %w", err, cancelErr)
			}
			return uuid.Nil, fmt.Errorf("payment failed, order %s cancelled: %w", order.ID, err)
		}

		if err := u.orderRepo.UpdateStatus(ctx, order.ID, nextStatus); err != nil {
			return uuid.Nil, fmt.Errorf("payment succeeded but failed to update order status: %w", err)
		}
	}

	return order.ID, nil
}

// executePayment は支払い方法に応じて決済サービスを呼び、注文の次のステータスを返す。
// 代引きは配達時に集金するため Paid にはせず Confirmed(注文確定)止まりにする。
func (u *orderManagementUsecase) executePayment(ctx context.Context, orderID uuid.UUID, input CreateOrderInput, totalAmount int64) (domain.OrderStatus, error) {
	if input.PaymentMethod == PaymentMethodCashOnDelivery {
		_, err := u.paymentClient.CreateCODPayment(ctx, client.CODPaymentInput{
			OrderID: orderID.String(),
			Amount:  totalAmount,
		})
		if err != nil {
			return "", err
		}
		return domain.OrderStatusConfirmed, nil
	}

	_, err := u.paymentClient.ProcessPayment(ctx, client.ProcessPaymentInput{
		OrderID:         orderID.String(),
		CustomerID:      input.CustomerID.String(),
		PaymentMethodID: input.PaymentMethodID,
		Amount:          totalAmount,
		Currency:        "jpy",
	})
	if err != nil {
		return "", err
	}
	return domain.OrderStatusPaid, nil
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
	// 支払い済みの決済があれば全額返金する(未決済なら決済サービスが NotFound を返し、返金はスキップされる)
	if u.paymentClient != nil {
		if err := u.paymentClient.RefundByOrder(ctx, orderID.String(), "order cancelled"); err != nil {
			return fmt.Errorf("failed to refund payment for order %s: %w", orderID, err)
		}
	}

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
