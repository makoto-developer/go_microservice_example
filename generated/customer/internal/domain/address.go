package domain

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID             uuid.UUID `db:"id" json:"id"`
	CustomerID     uuid.UUID `db:"customer_id" json:"customer_id"`
	AddressName    string    `db:"address_name" json:"address_name"`
	PostalCode     string    `db:"postal_code" json:"postal_code"`
	Prefecture     string    `db:"prefecture" json:"prefecture"`
	City           string    `db:"city" json:"city"`
	AddressLine1   string    `db:"address_line1" json:"address_line1"`
	AddressLine2   *string   `db:"address_line2" json:"address_line2,omitempty"`
	RecipientName  string    `db:"recipient_name" json:"recipient_name"`
	RecipientPhone string    `db:"recipient_phone" json:"recipient_phone"`
	IsDefault      bool      `db:"is_default" json:"is_default"`
	Deleted        bool      `db:"deleted" json:"deleted"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

func NewAddress(customerID uuid.UUID, addressName, postalCode, prefecture, city, addressLine1 string, addressLine2 *string, recipientName, recipientPhone string, isDefault bool) *Address {
	now := time.Now()
	return &Address{
		ID:             uuid.New(),
		CustomerID:     customerID,
		AddressName:    addressName,
		PostalCode:     postalCode,
		Prefecture:     prefecture,
		City:           city,
		AddressLine1:   addressLine1,
		AddressLine2:   addressLine2,
		RecipientName:  recipientName,
		RecipientPhone: recipientPhone,
		IsDefault:      isDefault,
		Deleted:        false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (a *Address) FullAddress() string {
	full := a.Prefecture + a.City + a.AddressLine1
	if a.AddressLine2 != nil {
		full += " " + *a.AddressLine2
	}
	return full
}
