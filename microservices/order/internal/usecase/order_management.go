package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/client"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/repository"
	"log"
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
	CustomerEmail   string // 通知メールの宛先(空なら通知スキップ)
	AddressID       uuid.UUID
	Items           []OrderItemInput
	ShippingFee     int64         // 0 の場合は shipping サービスに見積もりを依頼する
	ShippingMethod  string        // standard / express
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
	ListOrdersByCustomer(ctx context.Context, customerID uuid.UUID) ([]*domain.Order, error)
	Reorder(ctx context.Context, originalOrderID uuid.UUID) (uuid.UUID, error)
	CancelOrder(ctx context.Context, orderID uuid.UUID) error
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error
}

type orderManagementUsecase struct {
	orderRepo          repository.OrderRepository
	orderItemRepo      repository.OrderItemRepository
	paymentClient      client.PaymentClient
	shippingClient     client.ShippingClient     // nil の場合は出荷起票をスキップ
	notificationClient client.NotificationClient // nil の場合は通知をスキップ
	inventoryClient    client.InventoryClient    // nil の場合は在庫引当をスキップ
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

// NewOrderManagementUsecaseWithShipping は出荷起票(shipping サービス連携)込みで組み立てる。
func NewOrderManagementUsecaseWithShipping(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	paymentClient client.PaymentClient,
	shippingClient client.ShippingClient,
) OrderManagementUsecase {
	return &orderManagementUsecase{
		orderRepo:      orderRepo,
		orderItemRepo:  orderItemRepo,
		paymentClient:  paymentClient,
		shippingClient: shippingClient,
	}
}

// NewOrderManagementUsecaseFull は決済・出荷・通知の全連携込みで組み立てる。
func NewOrderManagementUsecaseFull(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	paymentClient client.PaymentClient,
	shippingClient client.ShippingClient,
	notificationClient client.NotificationClient,
	inventoryClient client.InventoryClient,
) OrderManagementUsecase {
	return &orderManagementUsecase{
		orderRepo:          orderRepo,
		orderItemRepo:      orderItemRepo,
		paymentClient:      paymentClient,
		shippingClient:     shippingClient,
		notificationClient: notificationClient,
		inventoryClient:    inventoryClient,
	}
}

func (u *orderManagementUsecase) CreateOrder(ctx context.Context, input CreateOrderInput) (uuid.UUID, error) {
	// 送料: 指定が無ければ shipping サービスに見積もりを依頼(失敗時は標準料金)
	if input.ShippingFee == 0 {
		input.ShippingFee = 500
		if u.shippingClient != nil {
			if fee, err := u.shippingClient.CalculateFee(ctx, input.ShippingMethod, 0); err == nil {
				input.ShippingFee = fee
			} else {
				log.Printf("shipping fee estimation failed (fallback to default): %v", err)
			}
		}
	}

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

	// 在庫の引当(inventory サービス)。在庫が無ければ決済前に注文を失敗させる
	if u.inventoryClient != nil {
		items := make([]client.StockItem, 0, len(input.Items))
		for _, item := range input.Items {
			items = append(items, client.StockItem{ProductID: item.ProductID.String(), Quantity: item.Quantity})
		}
		if err := u.inventoryClient.ReserveOrderStock(ctx, order.ID.String(), items); err != nil {
			if cancelErr := u.orderRepo.UpdateStatus(ctx, order.ID, domain.OrderStatusCancelled); cancelErr != nil {
				return uuid.Nil, fmt.Errorf("stock reservation failed (%v) and order cancellation also failed: %w", err, cancelErr)
			}
			return uuid.Nil, fmt.Errorf("stock reservation failed, order %s cancelled: %w", order.ID, err)
		}
	}

	// 決済サービスで支払いを実行する。失敗した場合は在庫引当を解放し、注文をキャンセルする。
	// paymentClient が未設定(nil)の場合は決済をスキップする(ローカル起動の後方互換)。
	if u.paymentClient != nil {
		nextStatus, err := u.executePayment(ctx, order.ID, input, totalAmount)
		if err != nil {
			u.releaseStock(ctx, order.ID) // 補償: 引当を解放
			if cancelErr := u.orderRepo.UpdateStatus(ctx, order.ID, domain.OrderStatusCancelled); cancelErr != nil {
				return uuid.Nil, fmt.Errorf("payment failed (%v) and order cancellation also failed: %w", err, cancelErr)
			}
			return uuid.Nil, fmt.Errorf("payment failed, order %s cancelled: %w", order.ID, err)
		}

		if err := u.orderRepo.UpdateStatus(ctx, order.ID, nextStatus); err != nil {
			return uuid.Nil, fmt.Errorf("payment succeeded but failed to update order status: %w", err)
		}
	}

	// 決済まで通ったので引当を確定する(失敗しても注文は成立、restock はバッチ想定)
	if u.inventoryClient != nil {
		if err := u.inventoryClient.ConfirmOrderStock(ctx, order.ID.String()); err != nil {
			log.Printf("stock confirmation failed for order %s: %v", order.ID, err)
		}
	}

	// 出荷の起票(shipping サービス)。失敗しても注文は成立させる(出荷はバックオフィスで再起票できる)
	if u.shippingClient != nil {
		if _, err := u.shippingClient.CreateShipment(ctx, client.CreateShipmentInput{
			OrderID:         order.ID.String(),
			CustomerID:      input.CustomerID.String(),
			ShippingAddress: input.AddressID.String(), // 住所解決は未実装のため ID をそのまま渡す
		}); err != nil {
			log.Printf("shipment creation failed for order %s (will be re-created manually): %v", order.ID, err)
		}
	}

	// 注文確認メール(通知サービス)。失敗しても注文は成立させる
	if u.notificationClient != nil && input.CustomerEmail != "" {
		if err := u.notificationClient.NotifyOrderConfirmed(ctx, client.OrderNotificationInput{
			CustomerID:    input.CustomerID.String(),
			CustomerEmail: input.CustomerEmail,
			OrderNumber:   order.OrderNumber,
			TotalAmount:   totalAmount,
		}); err != nil {
			log.Printf("order confirmation email failed for order %s: %v", order.ID, err)
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

// Reorder は過去の注文と同じ内容で新しい注文を作る(再注文)。
// 支払いは通常のカードフローに乗る(在庫引当・決済・出荷起票も通常どおり)。
func (u *orderManagementUsecase) Reorder(ctx context.Context, originalOrderID uuid.UUID) (uuid.UUID, error) {
	original, err := u.orderRepo.GetByID(ctx, originalOrderID)
	if err != nil || original == nil {
		return uuid.Nil, fmt.Errorf("original order not found: %w", err)
	}
	items, err := u.orderItemRepo.GetByOrderID(ctx, originalOrderID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to load original items: %w", err)
	}
	if len(items) == 0 {
		return uuid.Nil, fmt.Errorf("original order %s has no items", originalOrderID)
	}

	input := CreateOrderInput{
		CustomerID:      original.CustomerID,
		AddressID:       original.AddressID,
		ShippingFee:     original.ShippingFee,
		PaymentMethodID: "pm_reorder",
	}
	for _, item := range items {
		input.Items = append(input.Items, OrderItemInput{
			ProductID:   item.ProductID,
			VariationID: item.VariationID,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		})
	}
	return u.CreateOrder(ctx, input)
}

// releaseStock は在庫引当の解放(補償トランザクション)。失敗はログのみ(期限切れ解放バッチで回収される想定)。
func (u *orderManagementUsecase) releaseStock(ctx context.Context, orderID uuid.UUID) {
	if u.inventoryClient == nil {
		return
	}
	if err := u.inventoryClient.ReleaseOrderStock(ctx, orderID.String()); err != nil {
		log.Printf("stock release failed for order %s: %v", orderID, err)
	}
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

func (u *orderManagementUsecase) ListOrdersByCustomer(ctx context.Context, customerID uuid.UUID) ([]*domain.Order, error) {
	orders, err := u.orderRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	return orders, nil
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

	// 在庫引当を解放(best effort)
	u.releaseStock(ctx, orderID)

	// キャンセル確認メール(宛先が分かる場合のみ・best effort)
	if u.notificationClient != nil {
		if order, err := u.orderRepo.GetByID(ctx, orderID); err == nil && order != nil {
			if err := u.notificationClient.NotifyOrderCancelled(ctx, client.OrderNotificationInput{
				CustomerID:  order.CustomerID.String(),
				OrderNumber: order.OrderNumber,
				TotalAmount: order.TotalAmount,
			}); err != nil {
				log.Printf("order cancellation email failed for order %s: %v", orderID, err)
			}
		}
	}
	return nil
}

func (u *orderManagementUsecase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	if err := u.orderRepo.UpdateStatus(ctx, orderID, status); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}
