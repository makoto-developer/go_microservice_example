package domain

import (
	"time"

	"github.com/google/uuid"
)

// Shop represents Shop
type Shop struct {
	Id            uuid.UUID  `db:"id" json:"id"`
	OwnerId       uuid.UUID  `db:"owner_id" json:"owner_id"`
	Name          string     `db:"name" json:"name"`
	Description   string     `db:"description" json:"description"`
	LogoUrl       *string    `db:"logo_url" json:"logo_url,omitempty"`
	OwnerName     string     `db:"owner_name" json:"owner_name"`
	PhoneNumber   string     `db:"phone_number" json:"phone_number"`
	BusinessHours string     `db:"business_hours" json:"business_hours"`
	ReturnPolicy  string     `db:"return_policy" json:"return_policy"`
	Status        ShopStatus `db:"status" json:"status"`
	Published     bool       `db:"published" json:"published"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

// NewShop creates a new Shop instance
func NewShop() *Shop {
	return &Shop{}
}
