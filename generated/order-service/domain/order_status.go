package domain

// OrderStatus represents OrderStatus type
type OrderStatus string

const (
	OrderStatusPending OrderStatus = "PENDING"
	OrderStatusPaymentProcessing OrderStatus = "PAYMENT_PROCESSING"
	OrderStatusPaymentFailed OrderStatus = "PAYMENT_FAILED"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusPreparing OrderStatus = "PREPARING"
	OrderStatusShipped OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// OrderStatusValues returns all possible values
func OrderStatusValues() []OrderStatus {
	return []OrderStatus{
		OrderStatusPending,
		OrderStatusPaymentProcessing,
		OrderStatusPaymentFailed,
		OrderStatusConfirmed,
		OrderStatusPreparing,
		OrderStatusShipped,
		OrderStatusDelivered,
		OrderStatusCancelled,
	}
}

// IsValid checks if the value is valid
func (e OrderStatus) IsValid() bool {
	switch e {
	case OrderStatusPending:
	case OrderStatusPaymentProcessing:
	case OrderStatusPaymentFailed:
	case OrderStatusConfirmed:
	case OrderStatusPreparing:
	case OrderStatusShipped:
	case OrderStatusDelivered:
	case OrderStatusCancelled:
		return true
	}
	return false
}
