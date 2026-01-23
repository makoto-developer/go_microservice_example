package domain

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID    `db:"id" json:"id"`
	ShopID          uuid.UUID    `db:"shop_id" json:"shop_id"`
	CustomerID      uuid.UUID    `db:"customer_id" json:"customer_id"`
	OrderNumber     string       `db:"order_number" json:"order_number"`
	Status          OrderStatus  `db:"status" json:"status"`
	TotalAmount     float64      `db:"total_amount" json:"total_amount"`
	ShippingAddress string       `db:"shipping_address" json:"shipping_address"`
	PaymentMethod   string       `db:"payment_method" json:"payment_method"`
	TrackingNumber  *string      `db:"tracking_number" json:"tracking_number,omitempty"`
	Carrier         *Carrier     `db:"carrier" json:"carrier,omitempty"`
	CreatedAt       time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time    `db:"updated_at" json:"updated_at"`
}

func NewOrder(shopID, customerID uuid.UUID, orderNumber string, totalAmount float64, shippingAddr, paymentMethod string) *Order {
	now := time.Now()
	return &Order{
		ID:              uuid.New(),
		ShopID:          shopID,
		CustomerID:      customerID,
		OrderNumber:     orderNumber,
		Status:          OrderStatusReceived,
		TotalAmount:     totalAmount,
		ShippingAddress: shippingAddr,
		PaymentMethod:   paymentMethod,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (o *Order) CanUpdateStatus(newStatus OrderStatus) bool {
	transitions := map[OrderStatus][]OrderStatus{
		OrderStatusReceived:  {OrderStatusPreparing, OrderStatusCancelled},
		OrderStatusPreparing: {OrderStatusShipped, OrderStatusCancelled},
		OrderStatusShipped:   {OrderStatusDelivered},
		OrderStatusDelivered: {},
		OrderStatusCancelled: {},
	}

	allowedStatuses, ok := transitions[o.Status]
	if !ok {
		return false
	}

	for _, allowed := range allowedStatuses {
		if allowed == newStatus {
			return true
		}
	}
	return false
}
