package domain

// MetricType represents MetricType type
type MetricType string

const (
	MetricTypeTotalUsers MetricType = "TOTAL_USERS"
	MetricTypeActiveUsers MetricType = "ACTIVE_USERS"
	MetricTypeTotalShops MetricType = "TOTAL_SHOPS"
	MetricTypeActiveShops MetricType = "ACTIVE_SHOPS"
	MetricTypeTotalOrders MetricType = "TOTAL_ORDERS"
	MetricTypeTotalRevenue MetricType = "TOTAL_REVENUE"
	MetricTypeDailyOrders MetricType = "DAILY_ORDERS"
	MetricTypeDailyRevenue MetricType = "DAILY_REVENUE"
	MetricTypeNewUsers MetricType = "NEW_USERS"
	MetricTypeNewShops MetricType = "NEW_SHOPS"
)

// MetricTypeValues returns all possible values
func MetricTypeValues() []MetricType {
	return []MetricType{
		MetricTypeTotalUsers,
		MetricTypeActiveUsers,
		MetricTypeTotalShops,
		MetricTypeActiveShops,
		MetricTypeTotalOrders,
		MetricTypeTotalRevenue,
		MetricTypeDailyOrders,
		MetricTypeDailyRevenue,
		MetricTypeNewUsers,
		MetricTypeNewShops,
	}
}

// IsValid checks if the value is valid
func (e MetricType) IsValid() bool {
	switch e {
	case MetricTypeTotalUsers:
	case MetricTypeActiveUsers:
	case MetricTypeTotalShops:
	case MetricTypeActiveShops:
	case MetricTypeTotalOrders:
	case MetricTypeTotalRevenue:
	case MetricTypeDailyOrders:
	case MetricTypeDailyRevenue:
	case MetricTypeNewUsers:
	case MetricTypeNewShops:
		return true
	}
	return false
}
