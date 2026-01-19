package domain

// ShipmentStatus represents ShipmentStatus type
type ShipmentStatus string

const (
	ShipmentStatusPreparing ShipmentStatus = "PREPARING"
	ShipmentStatusAwaitingPickup ShipmentStatus = "AWAITING_PICKUP"
	ShipmentStatusInTransit ShipmentStatus = "IN_TRANSIT"
	ShipmentStatusAtDistributionCenter ShipmentStatus = "AT_DISTRIBUTION_CENTER"
	ShipmentStatusOutForDelivery ShipmentStatus = "OUT_FOR_DELIVERY"
	ShipmentStatusDelivered ShipmentStatus = "DELIVERED"
	ShipmentStatusDeliveryAttempted ShipmentStatus = "DELIVERY_ATTEMPTED"
	ShipmentStatusDeliveryFailed ShipmentStatus = "DELIVERY_FAILED"
)

// ShipmentStatusValues returns all possible values
func ShipmentStatusValues() []ShipmentStatus {
	return []ShipmentStatus{
		ShipmentStatusPreparing,
		ShipmentStatusAwaitingPickup,
		ShipmentStatusInTransit,
		ShipmentStatusAtDistributionCenter,
		ShipmentStatusOutForDelivery,
		ShipmentStatusDelivered,
		ShipmentStatusDeliveryAttempted,
		ShipmentStatusDeliveryFailed,
	}
}

// IsValid checks if the value is valid
func (e ShipmentStatus) IsValid() bool {
	switch e {
	case ShipmentStatusPreparing:
	case ShipmentStatusAwaitingPickup:
	case ShipmentStatusInTransit:
	case ShipmentStatusAtDistributionCenter:
	case ShipmentStatusOutForDelivery:
	case ShipmentStatusDelivered:
	case ShipmentStatusDeliveryAttempted:
	case ShipmentStatusDeliveryFailed:
		return true
	}
	return false
}
