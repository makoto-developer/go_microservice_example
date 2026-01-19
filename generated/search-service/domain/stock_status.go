package domain

// StockStatus represents StockStatus type
type StockStatus string

const (
	StockStatusInStock StockStatus = "IN_STOCK"
	StockStatusLowStock StockStatus = "LOW_STOCK"
	StockStatusOutOfStock StockStatus = "OUT_OF_STOCK"
)

// StockStatusValues returns all possible values
func StockStatusValues() []StockStatus {
	return []StockStatus{
		StockStatusInStock,
		StockStatusLowStock,
		StockStatusOutOfStock,
	}
}

// IsValid checks if the value is valid
func (e StockStatus) IsValid() bool {
	switch e {
	case StockStatusInStock:
	case StockStatusLowStock:
	case StockStatusOutOfStock:
		return true
	}
	return false
}
