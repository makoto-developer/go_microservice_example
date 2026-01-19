package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// Shipment represents Shipment
type Shipment struct {
	Id uuid.UUID `db:"id" json:"id"`
	OrderId uuid.UUID `db:"order_id" json:"order_id"`
	TrackingNumber *string `db:"tracking_number" json:"tracking_number,omitempty"`
	Carrier Carrier `db:"carrier" json:"carrier"`
	Status ShipmentStatus `db:"status" json:"status"`
	ShippingAddress string `db:"shipping_address" json:"shipping_address"`
	RecipientName string `db:"recipient_name" json:"recipient_name"`
	RecipientPhone string `db:"recipient_phone" json:"recipient_phone"`
	DeliveryDate *date `db:"delivery_date" json:"delivery_date,omitempty"`
	DeliveryTimeSlot *TimeSlot `db:"delivery_time_slot" json:"delivery_time_slot,omitempty"`
	DeliveryOption *DeliveryOption `db:"delivery_option" json:"delivery_option,omitempty"`
	ShippingFee decimal.Decimal `db:"shipping_fee" json:"shipping_fee"`
	DispatchedAt *time.Time `db:"dispatched_at" json:"dispatched_at,omitempty"`
	DeliveredAt *time.Time `db:"delivered_at" json:"delivered_at,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewShipment creates a new Shipment instance
func NewShipment() *Shipment {
	return &Shipment{}
}
