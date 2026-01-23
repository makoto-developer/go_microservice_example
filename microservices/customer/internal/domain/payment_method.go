package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod struct {
	ID                    uuid.UUID `db:"id" json:"id"`
	CustomerID            uuid.UUID `db:"customer_id" json:"customer_id"`
	StripePaymentMethodID string    `db:"stripe_payment_method_id" json:"stripe_payment_method_id"`
	CardLast4             string    `db:"card_last4" json:"card_last4"`
	CardBrand             string    `db:"card_brand" json:"card_brand"`
	CardExpMonth          int       `db:"card_exp_month" json:"card_exp_month"`
	CardExpYear           int       `db:"card_exp_year" json:"card_exp_year"`
	CardholderName        string    `db:"cardholder_name" json:"cardholder_name"`
	IsDefault             bool      `db:"is_default" json:"is_default"`
	CreatedAt             time.Time `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time `db:"updated_at" json:"updated_at"`
}

func NewPaymentMethod(customerID uuid.UUID, stripePaymentMethodID, cardLast4, cardBrand string, cardExpMonth, cardExpYear int, cardholderName string, isDefault bool) *PaymentMethod {
	now := time.Now()
	return &PaymentMethod{
		ID:                    uuid.New(),
		CustomerID:            customerID,
		StripePaymentMethodID: stripePaymentMethodID,
		CardLast4:             cardLast4,
		CardBrand:             cardBrand,
		CardExpMonth:          cardExpMonth,
		CardExpYear:           cardExpYear,
		CardholderName:        cardholderName,
		IsDefault:             isDefault,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
}

func (p *PaymentMethod) IsExpired() bool {
	now := time.Now()
	expYear := p.CardExpYear
	expMonth := p.CardExpMonth

	if now.Year() > expYear {
		return true
	}
	if now.Year() == expYear && int(now.Month()) > expMonth {
		return true
	}
	return false
}
