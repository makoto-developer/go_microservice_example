package domain

import (
	"time"

	"github.com/google/uuid"
)

// Customer represents a customer entity
type Customer struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	UserID          uuid.UUID  `db:"user_id" json:"user_id"`
	FirstName       string     `db:"first_name" json:"first_name"`
	LastName        string     `db:"last_name" json:"last_name"`
	PhoneNumber     string     `db:"phone_number" json:"phone_number"`
	BirthDate       *time.Time `db:"birth_date" json:"birth_date,omitempty"`
	Gender          *string    `db:"gender" json:"gender,omitempty"`
	ProfileImageURL *string    `db:"profile_image_url" json:"profile_image_url,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}
