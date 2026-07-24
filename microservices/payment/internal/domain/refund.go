package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusSucceeded RefundStatus = "succeeded"
	RefundStatusFailed    RefundStatus = "failed"
)

var (
	ErrPaymentNotFound     = errors.New("payment not found")
	ErrRefundNotFound      = errors.New("refund not found")
	ErrRefundNotAllowed    = errors.New("can only refund completed payments")
	ErrInvalidRefundAmount = errors.New("invalid refund amount")
	ErrOrderMismatch       = errors.New("payment does not belong to the order")
	ErrNotCODPayment       = errors.New("payment is not cash on delivery")
)

type Refund struct {
	ID            uuid.UUID    `db:"id" json:"id"`
	PaymentID     uuid.UUID    `db:"payment_id" json:"payment_id"`
	OrderID       uuid.UUID    `db:"order_id" json:"order_id"`
	Amount        int          `db:"amount" json:"amount"` // in cents
	Reason        string       `db:"reason" json:"reason"`
	Status        RefundStatus `db:"status" json:"status"`
	TransactionID string       `db:"transaction_id" json:"transaction_id"`
	CreatedAt     time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at" json:"updated_at"`
}

// NewRefund は完了済みの決済に対する返金を作る。金額は決済額を超えられない。
func NewRefund(payment *Payment, amount int, reason string) (*Refund, error) {
	if payment.Status != PaymentStatusCompleted {
		return nil, ErrRefundNotAllowed
	}
	if amount <= 0 || amount > payment.Amount {
		return nil, ErrInvalidRefundAmount
	}
	now := time.Now()
	return &Refund{
		ID:        uuid.New(),
		PaymentID: payment.ID,
		OrderID:   payment.OrderID,
		Amount:    amount,
		Reason:    reason,
		Status:    RefundStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r *Refund) Complete(transactionID string) {
	r.Status = RefundStatusSucceeded
	r.TransactionID = transactionID
	r.UpdatedAt = time.Now()
}
