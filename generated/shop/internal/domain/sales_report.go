package domain

import (
	"time"

	"github.com/google/uuid"
)

type SalesReport struct {
	ID                uuid.UUID `db:"id" json:"id"`
	ShopID            uuid.UUID `db:"shop_id" json:"shop_id"`
	Date              time.Time `db:"date" json:"date"`
	TotalSales        float64   `db:"total_sales" json:"total_sales"`
	OrderCount        int       `db:"order_count" json:"order_count"`
	AverageOrderValue float64   `db:"average_order_value" json:"average_order_value"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
}

func NewSalesReport(shopID uuid.UUID, date time.Time, totalSales float64, orderCount int) *SalesReport {
	avgOrderValue := 0.0
	if orderCount > 0 {
		avgOrderValue = totalSales / float64(orderCount)
	}

	return &SalesReport{
		ID:                uuid.New(),
		ShopID:            shopID,
		Date:              date,
		TotalSales:        totalSales,
		OrderCount:        orderCount,
		AverageOrderValue: avgOrderValue,
		CreatedAt:         time.Now(),
	}
}
