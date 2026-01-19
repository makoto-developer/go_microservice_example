package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

var (
	ErrInvalidOrderStatus = errors.New("invalid order status")
	ErrOrderAlreadyCancelled = errors.New("order already cancelled")
)

type Order struct {
	ID           uuid.UUID   `db:"id" json:"id"`
	CustomerID   uuid.UUID   `db:"customer_id" json:"customer_id"`
	OrderNumber  string      `db:"order_number" json:"order_number"`
	Status       OrderStatus `db:"status" json:"status"`
	TotalAmount  int         `db:"total_amount" json:"total_amount"` // in cents
	CreatedAt    time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	ID        uuid.UUID `db:"id" json:"id"`
	OrderID   uuid.UUID `db:"order_id" json:"order_id"`
	ProductID uuid.UUID `db:"product_id" json:"product_id"`
	ShopID    uuid.UUID `db:"shop_id" json:"shop_id"`
	Quantity  int       `db:"quantity" json:"quantity"`
	Price     int       `db:"price" json:"price"` // in cents
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewOrder(customerID uuid.UUID, orderNumber string, totalAmount int) *Order {
	now := time.Now()
	return &Order{
		ID:          uuid.New(),
		CustomerID:  customerID,
		OrderNumber: orderNumber,
		Status:      OrderStatusPending,
		TotalAmount: totalAmount,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func NewOrderItem(orderID, productID, shopID uuid.UUID, quantity, price int) *OrderItem {
	return &OrderItem{
		ID:        uuid.New(),
		OrderID:   orderID,
		ProductID: productID,
		ShopID:    shopID,
		Quantity:  quantity,
		Price:     price,
		CreatedAt: time.Now(),
	}
}

func (o *Order) CanCancel() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusConfirmed
}

func (o *Order) Cancel() error {
	if !o.CanCancel() {
		return ErrOrderAlreadyCancelled
	}
	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Confirm() error {
	if o.Status != OrderStatusPending {
		return ErrInvalidOrderStatus
	}
	o.Status = OrderStatusConfirmed
	o.UpdatedAt = time.Now()
	return nil
}
