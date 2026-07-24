package domain

import (
	"github.com/google/uuid"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID              uuid.UUID   `db:"id" json:"id"`
	CustomerID      uuid.UUID   `db:"customer_id" json:"customer_id"`
	OrderNumber     string      `db:"order_number" json:"order_number"`
	Status          OrderStatus `db:"status" json:"status"`
	TotalAmount     int64       `db:"total_amount" json:"total_amount"`
	ShippingFee     int64       `db:"shipping_fee" json:"shipping_fee"`
	AddressID       uuid.UUID   `db:"address_id" json:"address_id"`
	PaymentMethodID *uuid.UUID  `db:"payment_method_id" json:"payment_method_id,omitempty"`
	CreatedAt       time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	OrderID     uuid.UUID  `db:"order_id" json:"order_id"`
	ProductID   uuid.UUID  `db:"product_id" json:"product_id"`
	VariationID *uuid.UUID `db:"variation_id" json:"variation_id,omitempty"`
	Quantity    int        `db:"quantity" json:"quantity"`
	UnitPrice   int64      `db:"unit_price" json:"unit_price"`
	Subtotal    int64      `db:"subtotal" json:"subtotal"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}
