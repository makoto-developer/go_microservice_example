package domain

import (
	"time"

	"github.com/google/uuid"
)

// Shop represents a shop entity
type Shop struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	OwnerID        uuid.UUID  `db:"owner_id" json:"owner_id"`
	Name           string     `db:"name" json:"name"`
	Description    string     `db:"description" json:"description"`
	LogoImageURL   string     `db:"logo_image_url" json:"logo_image_url"`
	OwnerName      string     `db:"owner_name" json:"owner_name"`
	OwnerPhone     string     `db:"owner_phone" json:"owner_phone"`
	BusinessHours  string     `db:"business_hours" json:"business_hours"`
	ReturnPolicy   string     `db:"return_policy" json:"return_policy"`
	Status         ShopStatus `db:"status" json:"status"`
	IsPublic       bool       `db:"is_public" json:"is_public"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	ApprovedAt     *time.Time `db:"approved_at" json:"approved_at,omitempty"`
	ApprovedBy     *uuid.UUID `db:"approved_by" json:"approved_by,omitempty"`
}

// Category represents a shop category
type Category struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
