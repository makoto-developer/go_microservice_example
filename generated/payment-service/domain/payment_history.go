package domain

import (
	"github.com/google/uuid"
	"time"
)

// PaymentHistory represents PaymentHistory
type PaymentHistory struct {
	Id uuid.UUID `db:"id" json:"id"`
	PaymentId uuid.UUID `db:"payment_id" json:"payment_id"`
	OldStatus *PaymentStatus `db:"old_status" json:"old_status,omitempty"`
	NewStatus PaymentStatus `db:"new_status" json:"new_status"`
	ChangedBy string `db:"changed_by" json:"changed_by"`
	ChangeReason *string `db:"change_reason" json:"change_reason,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewPaymentHistory creates a new PaymentHistory instance
func NewPaymentHistory() *PaymentHistory {
	return &PaymentHistory{}
}
