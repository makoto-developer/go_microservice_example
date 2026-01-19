package domain

import (
	"time"
	"github.com/google/uuid"
)

// Address represents Address
type Address struct {
	Id uuid.UUID `db:"id" json:"id"`
	CustomerId uuid.UUID `db:"customer_id" json:"customer_id"`
	AddressName string `db:"address_name" json:"address_name"`
	PostalCode string `db:"postal_code" json:"postal_code"`
	Prefecture string `db:"prefecture" json:"prefecture"`
	City string `db:"city" json:"city"`
	AddressLine1 string `db:"address_line1" json:"address_line1"`
	AddressLine2 *string `db:"address_line2" json:"address_line2,omitempty"`
	RecipientName string `db:"recipient_name" json:"recipient_name"`
	RecipientPhone string `db:"recipient_phone" json:"recipient_phone"`
	IsDefault bool `db:"is_default" json:"is_default"`
	Deleted bool `db:"deleted" json:"deleted"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewAddress creates a new Address instance
func NewAddress() *Address {
	return &Address{}
}
