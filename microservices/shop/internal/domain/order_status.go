package domain

type OrderStatus string

const (
	OrderStatusReceived  OrderStatus = "RECEIVED"
	OrderStatusPreparing OrderStatus = "PREPARING"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

func (s OrderStatus) String() string {
	return string(s)
}

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusReceived, OrderStatusPreparing, OrderStatusShipped, OrderStatusDelivered, OrderStatusCancelled:
		return true
	}
	return false
}
