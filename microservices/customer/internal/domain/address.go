package domain

import (
	"time"

	"github.com/google/uuid"
)

// Address represents a shipping address
type Address struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	CustomerID     uuid.UUID  `db:"customer_id" json:"customer_id"`
	Label          string     `db:"label" json:"label"`
	PostalCode     string     `db:"postal_code" json:"postal_code"`
	Prefecture     string     `db:"prefecture" json:"prefecture"`
	City           string     `db:"city" json:"city"`
	AddressLine1   string     `db:"address_line1" json:"address_line1"`
	AddressLine2   *string    `db:"address_line2" json:"address_line2,omitempty"`
	RecipientName  string     `db:"recipient_name" json:"recipient_name"`
	RecipientPhone string     `db:"recipient_phone" json:"recipient_phone"`
	IsDefault      bool       `db:"is_default" json:"is_default"`
	IsDeleted      bool       `db:"is_deleted" json:"is_deleted"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// NewAddress creates a new Address
func NewAddress(customerID uuid.UUID, label, postalCode, prefecture, city, addressLine1, recipientName, recipientPhone string, addressLine2 *string) *Address {
	return &Address{
		ID:             uuid.New(),
		CustomerID:     customerID,
		Label:          label,
		PostalCode:     postalCode,
		Prefecture:     prefecture,
		City:           city,
		AddressLine1:   addressLine1,
		AddressLine2:   addressLine2,
		RecipientName:  recipientName,
		RecipientPhone: recipientPhone,
		IsDefault:      false,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
