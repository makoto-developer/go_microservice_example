package domain

import (
	"github.com/google/uuid"
	"time"
)

// OrderStatusHistory represents OrderStatusHistory
type OrderStatusHistory struct {
	Id uuid.UUID `db:"id" json:"id"`
	OrderId uuid.UUID `db:"order_id" json:"order_id"`
	OldStatus *OrderStatus `db:"old_status" json:"old_status,omitempty"`
	NewStatus OrderStatus `db:"new_status" json:"new_status"`
	ChangedBy string `db:"changed_by" json:"changed_by"`
	ChangeReason *string `db:"change_reason" json:"change_reason,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewOrderStatusHistory creates a new OrderStatusHistory instance
func NewOrderStatusHistory() *OrderStatusHistory {
	return &OrderStatusHistory{}
}
