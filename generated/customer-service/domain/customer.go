package domain

import (
	"github.com/google/uuid"
	"time"
)

// Customer represents Customer
type Customer struct {
	Id uuid.UUID `db:"id" json:"id"`
	UserId uuid.UUID `db:"user_id" json:"user_id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName string `db:"last_name" json:"last_name"`
	Phone string `db:"phone" json:"phone"`
	BirthDate *date `db:"birth_date" json:"birth_date,omitempty"`
	Gender *Gender `db:"gender" json:"gender,omitempty"`
	ProfileImageUrl *string `db:"profile_image_url" json:"profile_image_url,omitempty"`
	ProfileThumbnail100Url *string `db:"profile_thumbnail_100_url" json:"profile_thumbnail_100_url,omitempty"`
	ProfileThumbnail200Url *string `db:"profile_thumbnail_200_url" json:"profile_thumbnail_200_url,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewCustomer creates a new Customer instance
func NewCustomer() *Customer {
	return &Customer{}
}
