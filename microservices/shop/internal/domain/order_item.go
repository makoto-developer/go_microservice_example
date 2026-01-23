package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID          uuid.UUID `db:"id" json:"id"`
	OrderID     uuid.UUID `db:"order_id" json:"order_id"`
	ProductID   uuid.UUID `db:"product_id" json:"product_id"`
	ProductName string    `db:"product_name" json:"product_name"`
	Quantity    int       `db:"quantity" json:"quantity"`
	UnitPrice   float64   `db:"unit_price" json:"unit_price"`
	Subtotal    float64   `db:"subtotal" json:"subtotal"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func NewOrderItem(orderID, productID uuid.UUID, productName string, quantity int, unitPrice float64) *OrderItem {
	return &OrderItem{
		ID:          uuid.New(),
		OrderID:     orderID,
		ProductID:   productID,
		ProductName: productName,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		Subtotal:    float64(quantity) * unitPrice,
		CreatedAt:   time.Now(),
	}
}
