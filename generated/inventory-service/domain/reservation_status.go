package domain

// ReservationStatus represents ReservationStatus type
type ReservationStatus string

const (
	ReservationStatusReserved ReservationStatus = "RESERVED"
	ReservationStatusConfirmed ReservationStatus = "CONFIRMED"
	ReservationStatusReleased ReservationStatus = "RELEASED"
	ReservationStatusExpired ReservationStatus = "EXPIRED"
)

// ReservationStatusValues returns all possible values
func ReservationStatusValues() []ReservationStatus {
	return []ReservationStatus{
		ReservationStatusReserved,
		ReservationStatusConfirmed,
		ReservationStatusReleased,
		ReservationStatusExpired,
	}
}

// IsValid checks if the value is valid
func (e ReservationStatus) IsValid() bool {
	switch e {
	case ReservationStatusReserved:
	case ReservationStatusConfirmed:
	case ReservationStatusReleased:
	case ReservationStatusExpired:
		return true
	}
	return false
}
