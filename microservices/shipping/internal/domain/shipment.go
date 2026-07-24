package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ShipmentStatus string

const (
	ShipmentStatusPending   ShipmentStatus = "pending"
	ShipmentStatusPreparing ShipmentStatus = "preparing"
	ShipmentStatusShipped   ShipmentStatus = "shipped"
	ShipmentStatusInTransit ShipmentStatus = "in_transit"
	ShipmentStatusDelivered ShipmentStatus = "delivered"
	ShipmentStatusFailed    ShipmentStatus = "failed"
)

var (
	ErrInvalidStatusTransition  = errors.New("invalid status transition")
	ErrShipmentAlreadyDelivered = errors.New("shipment already delivered")
)

type Shipment struct {
	ID                uuid.UUID      `db:"id" json:"id"`
	OrderID           uuid.UUID      `db:"order_id" json:"order_id"`
	CustomerID        uuid.UUID      `db:"customer_id" json:"customer_id"`
	Status            ShipmentStatus `db:"status" json:"status"`
	TrackingNumber    string         `db:"tracking_number" json:"tracking_number"`
	Carrier           string         `db:"carrier" json:"carrier"`
	ShippingAddress   string         `db:"shipping_address" json:"shipping_address"`
	EstimatedDelivery time.Time      `db:"estimated_delivery" json:"estimated_delivery"`
	ActualDelivery    *time.Time     `db:"actual_delivery" json:"actual_delivery,omitempty"`
	CreatedAt         time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at" json:"updated_at"`
}

func NewShipment(orderID, customerID uuid.UUID, shippingAddress, carrier string) *Shipment {
	now := time.Now()
	return &Shipment{
		ID:                uuid.New(),
		OrderID:           orderID,
		CustomerID:        customerID,
		Status:            ShipmentStatusPending,
		Carrier:           carrier,
		ShippingAddress:   shippingAddress,
		EstimatedDelivery: now.Add(72 * time.Hour), // 3 days default
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func (s *Shipment) Ship(trackingNumber string) error {
	if s.Status != ShipmentStatusPreparing {
		return ErrInvalidStatusTransition
	}
	s.Status = ShipmentStatusShipped
	s.TrackingNumber = trackingNumber
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Shipment) MarkInTransit() error {
	if s.Status != ShipmentStatusShipped {
		return ErrInvalidStatusTransition
	}
	s.Status = ShipmentStatusInTransit
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Shipment) Deliver() error {
	if s.Status != ShipmentStatusInTransit {
		return ErrInvalidStatusTransition
	}
	now := time.Now()
	s.Status = ShipmentStatusDelivered
	s.ActualDelivery = &now
	s.UpdatedAt = now
	return nil
}

func (s *Shipment) CanUpdate() bool {
	return s.Status != ShipmentStatusDelivered && s.Status != ShipmentStatusFailed
}
