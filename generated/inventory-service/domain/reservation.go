package domain

import (
	"github.com/google/uuid"
	"time"
)

// Reservation represents Reservation
type Reservation struct {
	Id uuid.UUID `db:"id" json:"id"`
	InventoryId uuid.UUID `db:"inventory_id" json:"inventory_id"`
	OrderId uuid.UUID `db:"order_id" json:"order_id"`
	Quantity int `db:"quantity" json:"quantity"`
	Status ReservationStatus `db:"status" json:"status"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	ConfirmedAt *time.Time `db:"confirmed_at" json:"confirmed_at,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewReservation creates a new Reservation instance
func NewReservation() *Reservation {
	return &Reservation{}
}
