package domain

import (
	"github.com/google/uuid"
	"time"
)

// PaymentMethod represents PaymentMethod
type PaymentMethod struct {
	Id uuid.UUID `db:"id" json:"id"`
	CustomerId uuid.UUID `db:"customer_id" json:"customer_id"`
	StripePaymentMethodId string `db:"stripe_payment_method_id" json:"stripe_payment_method_id"`
	CardLast4 string `db:"card_last4" json:"card_last4"`
	CardBrand string `db:"card_brand" json:"card_brand"`
	CardExpMonth int `db:"card_exp_month" json:"card_exp_month"`
	CardExpYear int `db:"card_exp_year" json:"card_exp_year"`
	CardholderName string `db:"cardholder_name" json:"cardholder_name"`
	IsDefault bool `db:"is_default" json:"is_default"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewPaymentMethod creates a new PaymentMethod instance
func NewPaymentMethod() *PaymentMethod {
	return &PaymentMethod{}
}
