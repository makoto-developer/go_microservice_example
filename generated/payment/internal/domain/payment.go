package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string
type PaymentMethod string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodDebitCard  PaymentMethod = "debit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
)

var (
	ErrPaymentAlreadyProcessed = errors.New("payment already processed")
	ErrInvalidAmount = errors.New("invalid payment amount")
)

type Payment struct {
	ID            uuid.UUID     `db:"id" json:"id"`
	OrderID       uuid.UUID     `db:"order_id" json:"order_id"`
	Amount        int           `db:"amount" json:"amount"` // in cents
	Status        PaymentStatus `db:"status" json:"status"`
	PaymentMethod PaymentMethod `db:"payment_method" json:"payment_method"`
	TransactionID string        `db:"transaction_id" json:"transaction_id"`
	CreatedAt     time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at" json:"updated_at"`
}

func NewPayment(orderID uuid.UUID, amount int, method PaymentMethod) *Payment {
	now := time.Now()
	return &Payment{
		ID:            uuid.New(),
		OrderID:       orderID,
		Amount:        amount,
		Status:        PaymentStatusPending,
		PaymentMethod: method,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (p *Payment) CanProcess() bool {
	return p.Status == PaymentStatusPending
}

func (p *Payment) Complete(transactionID string) error {
	if !p.CanProcess() {
		return ErrPaymentAlreadyProcessed
	}
	p.Status = PaymentStatusCompleted
	p.TransactionID = transactionID
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Payment) Fail() error {
	if !p.CanProcess() {
		return ErrPaymentAlreadyProcessed
	}
	p.Status = PaymentStatusFailed
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Payment) Refund() error {
	if p.Status != PaymentStatusCompleted {
		return errors.New("can only refund completed payments")
	}
	p.Status = PaymentStatusRefunded
	p.UpdatedAt = time.Now()
	return nil
}
