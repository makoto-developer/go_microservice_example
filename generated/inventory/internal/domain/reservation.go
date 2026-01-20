package domain

import (
	"time"
	"github.com/google/uuid"
)

type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "PENDING"
	ReservationStatusConfirmed ReservationStatus = "CONFIRMED"
	ReservationStatusReleased  ReservationStatus = "RELEASED"
	ReservationStatusExpired   ReservationStatus = "EXPIRED"
)

type Reservation struct {
	ID          uuid.UUID         `db:"id" json:"id"`
	InventoryID uuid.UUID         `db:"inventory_id" json:"inventory_id"`
	OrderID     uuid.UUID         `db:"order_id" json:"order_id"`
	Quantity    int               `db:"quantity" json:"quantity"`
	Status      ReservationStatus `db:"status" json:"status"`
	ExpiresAt   time.Time         `db:"expires_at" json:"expires_at"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at"`
}
