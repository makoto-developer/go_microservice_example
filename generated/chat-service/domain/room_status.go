package domain

// RoomStatus represents RoomStatus type
type RoomStatus string

const (
	RoomStatusActive RoomStatus = "ACTIVE"
	RoomStatusResolved RoomStatus = "RESOLVED"
	RoomStatusClosed RoomStatus = "CLOSED"
)

// RoomStatusValues returns all possible values
func RoomStatusValues() []RoomStatus {
	return []RoomStatus{
		RoomStatusActive,
		RoomStatusResolved,
		RoomStatusClosed,
	}
}

// IsValid checks if the value is valid
func (e RoomStatus) IsValid() bool {
	switch e {
	case RoomStatusActive:
	case RoomStatusResolved:
	case RoomStatusClosed:
		return true
	}
	return false
}
