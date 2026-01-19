package domain

// InventoryStatus represents InventoryStatus type
type InventoryStatus string

const (
	InventoryStatusInStock InventoryStatus = "IN_STOCK"
	InventoryStatusLowStock InventoryStatus = "LOW_STOCK"
	InventoryStatusOutOfStock InventoryStatus = "OUT_OF_STOCK"
	InventoryStatusReserved InventoryStatus = "RESERVED"
)

// InventoryStatusValues returns all possible values
func InventoryStatusValues() []InventoryStatus {
	return []InventoryStatus{
		InventoryStatusInStock,
		InventoryStatusLowStock,
		InventoryStatusOutOfStock,
		InventoryStatusReserved,
	}
}

// IsValid checks if the value is valid
func (e InventoryStatus) IsValid() bool {
	switch e {
	case InventoryStatusInStock:
	case InventoryStatusLowStock:
	case InventoryStatusOutOfStock:
	case InventoryStatusReserved:
		return true
	}
	return false
}
