package domain

import (
	"github.com/google/uuid"
	"time"
)

// OrderCancellation represents OrderCancellation
type OrderCancellation struct {
	Id uuid.UUID `db:"id" json:"id"`
	OrderId uuid.UUID `db:"order_id" json:"order_id"`
	CancelledBy string `db:"cancelled_by" json:"cancelled_by"`
	CancelReason CancelReason `db:"cancel_reason" json:"cancel_reason"`
	CancelNote *text `db:"cancel_note" json:"cancel_note,omitempty"`
	RefundId *uuid.UUID `db:"refund_id" json:"refund_id,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewOrderCancellation creates a new OrderCancellation instance
func NewOrderCancellation() *OrderCancellation {
	return &OrderCancellation{}
}
