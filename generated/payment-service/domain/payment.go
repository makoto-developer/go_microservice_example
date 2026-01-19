package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// Payment represents Payment
type Payment struct {
	Id uuid.UUID `db:"id" json:"id"`
	OrderId uuid.UUID `db:"order_id" json:"order_id"`
	PaymentMethod PaymentMethodType `db:"payment_method" json:"payment_method"`
	Amount decimal.Decimal `db:"amount" json:"amount"`
	Currency string `db:"currency" json:"currency"`
	Status PaymentStatus `db:"status" json:"status"`
	StripePaymentIntentId *string `db:"stripe_payment_intent_id" json:"stripe_payment_intent_id,omitempty"`
	StripeClientSecret *string `db:"stripe_client_secret" json:"stripe_client_secret,omitempty"`
	RequiresAuthentication bool `db:"requires_authentication" json:"requires_authentication"`
	AuthenticationUrl *string `db:"authentication_url" json:"authentication_url,omitempty"`
	ErrorCode *string `db:"error_code" json:"error_code,omitempty"`
	ErrorMessage *text `db:"error_message" json:"error_message,omitempty"`
	CodFee *decimal.Decimal `db:"cod_fee" json:"cod_fee,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewPayment creates a new Payment instance
func NewPayment() *Payment {
	return &Payment{}
}
