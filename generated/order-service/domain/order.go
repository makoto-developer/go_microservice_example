package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// Order represents Order
type Order struct {
	Id uuid.UUID `db:"id" json:"id"`
	OrderNumber string `db:"order_number" json:"order_number"`
	CustomerId uuid.UUID `db:"customer_id" json:"customer_id"`
	CustomerEmail string `db:"customer_email" json:"customer_email"`
	Status OrderStatus `db:"status" json:"status"`
	TotalAmount decimal.Decimal `db:"total_amount" json:"total_amount"`
	Subtotal decimal.Decimal `db:"subtotal" json:"subtotal"`
	ShippingFee decimal.Decimal `db:"shipping_fee" json:"shipping_fee"`
	PaymentMethod PaymentMethod `db:"payment_method" json:"payment_method"`
	PaymentId *uuid.UUID `db:"payment_id" json:"payment_id,omitempty"`
	ShippingMethod string `db:"shipping_method" json:"shipping_method"`
	ShippingAddress string `db:"shipping_address" json:"shipping_address"`
	RecipientName string `db:"recipient_name" json:"recipient_name"`
	RecipientPhone string `db:"recipient_phone" json:"recipient_phone"`
	TrackingNumber *string `db:"tracking_number" json:"tracking_number,omitempty"`
	Carrier *string `db:"carrier" json:"carrier,omitempty"`
	Notes *text `db:"notes" json:"notes,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewOrder creates a new Order instance
func NewOrder() *Order {
	return &Order{}
}
