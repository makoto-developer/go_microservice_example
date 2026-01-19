package domain

// OrderStatus represents OrderStatus type
type OrderStatus string

const (
	OrderStatusReceived  OrderStatus = "RECEIVED"
	OrderStatusPreparing OrderStatus = "PREPARING"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// OrderStatusValues returns all possible values
func OrderStatusValues() []OrderStatus {
	return []OrderStatus{
		OrderStatusReceived,
		OrderStatusPreparing,
		OrderStatusShipped,
		OrderStatusDelivered,
		OrderStatusCancelled,
	}
}

// IsValid checks if the value is valid
func (e OrderStatus) IsValid() bool {
	switch e {
	case OrderStatusReceived:
	case OrderStatusPreparing:
	case OrderStatusShipped:
	case OrderStatusDelivered:
	case OrderStatusCancelled:
		return true
	}
	return false
}
