package domain

import (
	"time"
	"github.com/google/uuid"
)

// CarrierTracking represents CarrierTracking
type CarrierTracking struct {
	Id uuid.UUID `db:"id" json:"id"`
	ShipmentId uuid.UUID `db:"shipment_id" json:"shipment_id"`
	TrackingNumber string `db:"tracking_number" json:"tracking_number"`
	Carrier Carrier `db:"carrier" json:"carrier"`
	LastUpdatedAt time.Time `db:"last_updated_at" json:"last_updated_at"`
	TrackingData map[string]interface{} `db:"tracking_data" json:"tracking_data"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewCarrierTracking creates a new CarrierTracking instance
func NewCarrierTracking() *CarrierTracking {
	return &CarrierTracking{}
}
