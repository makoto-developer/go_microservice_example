package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// SalesReport represents SalesReport
type SalesReport struct {
	Id                uuid.UUID       `db:"id" json:"id"`
	ShopId            uuid.UUID       `db:"shop_id" json:"shop_id"`
	Date              time.Time       `db:"date" json:"date"`
	TotalSales        decimal.Decimal `db:"total_sales" json:"total_sales"`
	OrderCount        int             `db:"order_count" json:"order_count"`
	AverageOrderValue decimal.Decimal `db:"average_order_value" json:"average_order_value"`
	CreatedAt         time.Time       `db:"created_at" json:"created_at"`
}

// NewSalesReport creates a new SalesReport instance
func NewSalesReport() *SalesReport {
	return &SalesReport{}
}
