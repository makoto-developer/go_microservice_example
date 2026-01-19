package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// OrderItem represents OrderItem
type OrderItem struct {
	Id uuid.UUID `db:"id" json:"id"`
	OrderId uuid.UUID `db:"order_id" json:"order_id"`
	ShopId uuid.UUID `db:"shop_id" json:"shop_id"`
	ProductId uuid.UUID `db:"product_id" json:"product_id"`
	VariationId *uuid.UUID `db:"variation_id" json:"variation_id,omitempty"`
	ProductName string `db:"product_name" json:"product_name"`
	ProductSnapshot map[string]interface{} `db:"product_snapshot" json:"product_snapshot"`
	UnitPrice decimal.Decimal `db:"unit_price" json:"unit_price"`
	Quantity int `db:"quantity" json:"quantity"`
	Subtotal decimal.Decimal `db:"subtotal" json:"subtotal"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewOrderItem creates a new OrderItem instance
func NewOrderItem() *OrderItem {
	return &OrderItem{}
}
