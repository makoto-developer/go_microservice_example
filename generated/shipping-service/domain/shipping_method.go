package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// ShippingMethod represents ShippingMethod
type ShippingMethod struct {
	Id uuid.UUID `db:"id" json:"id"`
	ShopId uuid.UUID `db:"shop_id" json:"shop_id"`
	Name string `db:"name" json:"name"`
	Carrier Carrier `db:"carrier" json:"carrier"`
	BaseFee decimal.Decimal `db:"base_fee" json:"base_fee"`
	WeightLimitKg *decimal.Decimal `db:"weight_limit_kg" json:"weight_limit_kg,omitempty"`
	SizeLimitCm *int `db:"size_limit_cm" json:"size_limit_cm,omitempty"`
	IsActive bool `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewShippingMethod creates a new ShippingMethod instance
func NewShippingMethod() *ShippingMethod {
	return &ShippingMethod{}
}
