package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// OrderItem represents OrderItem
type OrderItem struct {
	Id          uuid.UUID       `db:"id" json:"id"`
	OrderId     uuid.UUID       `db:"order_id" json:"order_id"`
	ProductId   uuid.UUID       `db:"product_id" json:"product_id"`
	ProductName string          `db:"product_name" json:"product_name"`
	Quantity    int             `db:"quantity" json:"quantity"`
	UnitPrice   decimal.Decimal `db:"unit_price" json:"unit_price"`
	Subtotal    decimal.Decimal `db:"subtotal" json:"subtotal"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
}

// NewOrderItem creates a new OrderItem instance
func NewOrderItem() *OrderItem {
	return &OrderItem{}
}
