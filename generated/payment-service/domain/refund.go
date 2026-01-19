package domain

import (
	"time"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Refund represents Refund
type Refund struct {
	Id uuid.UUID `db:"id" json:"id"`
	PaymentId uuid.UUID `db:"payment_id" json:"payment_id"`
	OrderId uuid.UUID `db:"order_id" json:"order_id"`
	Amount decimal.Decimal `db:"amount" json:"amount"`
	Reason string `db:"reason" json:"reason"`
	Status RefundStatus `db:"status" json:"status"`
	StripeRefundId *string `db:"stripe_refund_id" json:"stripe_refund_id,omitempty"`
	ErrorMessage *text `db:"error_message" json:"error_message,omitempty"`
	RefundedAt *time.Time `db:"refunded_at" json:"refunded_at,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewRefund creates a new Refund instance
func NewRefund() *Refund {
	return &Refund{}
}
