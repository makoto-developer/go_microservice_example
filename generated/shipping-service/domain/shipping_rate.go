package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// ShippingRate represents ShippingRate
type ShippingRate struct {
	Id uuid.UUID `db:"id" json:"id"`
	ShippingMethodId uuid.UUID `db:"shipping_method_id" json:"shipping_method_id"`
	Prefecture string `db:"prefecture" json:"prefecture"`
	BaseRate decimal.Decimal `db:"base_rate" json:"base_rate"`
	WeightRatePerKg decimal.Decimal `db:"weight_rate_per_kg" json:"weight_rate_per_kg"`
	SizeRatePerCm decimal.Decimal `db:"size_rate_per_cm" json:"size_rate_per_cm"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewShippingRate creates a new ShippingRate instance
func NewShippingRate() *ShippingRate {
	return &ShippingRate{}
}
