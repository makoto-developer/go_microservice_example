package domain

import (
	"github.com/google/uuid"
	"time"
)

// ShipmentHistory represents ShipmentHistory
type ShipmentHistory struct {
	Id uuid.UUID `db:"id" json:"id"`
	ShipmentId uuid.UUID `db:"shipment_id" json:"shipment_id"`
	OldStatus *ShipmentStatus `db:"old_status" json:"old_status,omitempty"`
	NewStatus ShipmentStatus `db:"new_status" json:"new_status"`
	Location *string `db:"location" json:"location,omitempty"`
	Description *text `db:"description" json:"description,omitempty"`
	UpdatedBy string `db:"updated_by" json:"updated_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewShipmentHistory creates a new ShipmentHistory instance
func NewShipmentHistory() *ShipmentHistory {
	return &ShipmentHistory{}
}
