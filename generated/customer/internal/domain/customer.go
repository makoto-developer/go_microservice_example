package domain

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID                     uuid.UUID  `db:"id" json:"id"`
	UserID                 uuid.UUID  `db:"user_id" json:"user_id"`
	FirstName              string     `db:"first_name" json:"first_name"`
	LastName               string     `db:"last_name" json:"last_name"`
	Phone                  string     `db:"phone" json:"phone"`
	BirthDate              *time.Time `db:"birth_date" json:"birth_date,omitempty"`
	Gender                 *Gender    `db:"gender" json:"gender,omitempty"`
	ProfileImageURL        *string    `db:"profile_image_url" json:"profile_image_url,omitempty"`
	ProfileThumbnail100URL *string    `db:"profile_thumbnail_100_url" json:"profile_thumbnail_100_url,omitempty"`
	ProfileThumbnail200URL *string    `db:"profile_thumbnail_200_url" json:"profile_thumbnail_200_url,omitempty"`
	CreatedAt              time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time  `db:"updated_at" json:"updated_at"`
}

func NewCustomer(userID uuid.UUID, firstName, lastName, phone string) *Customer {
	now := time.Now()
	return &Customer{
		ID:        uuid.New(),
		UserID:    userID,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (c *Customer) FullName() string {
	return c.FirstName + " " + c.LastName
}
