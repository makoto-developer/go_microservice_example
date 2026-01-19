package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Order represents Order
type Order struct {
	Id              uuid.UUID       `db:"id" json:"id"`
	ShopId          uuid.UUID       `db:"shop_id" json:"shop_id"`
	CustomerId      uuid.UUID       `db:"customer_id" json:"customer_id"`
	OrderNumber     string          `db:"order_number" json:"order_number"`
	Status          OrderStatus     `db:"status" json:"status"`
	TotalAmount     decimal.Decimal `db:"total_amount" json:"total_amount"`
	ShippingAddress string          `db:"shipping_address" json:"shipping_address"`
	PaymentMethod   string          `db:"payment_method" json:"payment_method"`
	TrackingNumber  *string         `db:"tracking_number" json:"tracking_number,omitempty"`
	Carrier         *Carrier        `db:"carrier" json:"carrier,omitempty"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
}

// NewOrder creates a new Order instance
func NewOrder() *Order {
	return &Order{}
}
