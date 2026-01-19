package domain

import (
	"time"

	"github.com/google/uuid"
)

type Shop struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	OwnerID       uuid.UUID  `db:"owner_id" json:"owner_id"`
	Name          string     `db:"name" json:"name"`
	Description   string     `db:"description" json:"description"`
	LogoURL       *string    `db:"logo_url" json:"logo_url,omitempty"`
	OwnerName     string     `db:"owner_name" json:"owner_name"`
	PhoneNumber   string     `db:"phone_number" json:"phone_number"`
	BusinessHours string     `db:"business_hours" json:"business_hours"`
	ReturnPolicy  string     `db:"return_policy" json:"return_policy"`
	Status        ShopStatus `db:"status" json:"status"`
	Published     bool       `db:"published" json:"published"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

func NewShop(ownerID uuid.UUID, name, description string, logoURL *string, ownerName, phoneNumber, businessHours, returnPolicy string) *Shop {
	now := time.Now()
	return &Shop{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		Name:          name,
		Description:   description,
		LogoURL:       logoURL,
		OwnerName:     ownerName,
		PhoneNumber:   phoneNumber,
		BusinessHours: businessHours,
		ReturnPolicy:  returnPolicy,
		Status:        ShopStatusPendingApproval,
		Published:     false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (s *Shop) IsApproved() bool {
	return s.Status == ShopStatusApproved
}

func (s *Shop) CanPublish() bool {
	return s.Status == ShopStatusApproved
}
