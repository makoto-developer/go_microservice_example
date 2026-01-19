package domain

import (
	"github.com/shopspring/decimal"
	"time"
	"github.com/google/uuid"
)

// DashboardMetrics represents DashboardMetrics
type DashboardMetrics struct {
	Id uuid.UUID `db:"id" json:"id"`
	MetricType MetricType `db:"metric_type" json:"metric_type"`
	MetricValue decimal.Decimal `db:"metric_value" json:"metric_value"`
	MetricDate date `db:"metric_date" json:"metric_date"`
	Metadata *map[string]interface{} `db:"metadata" json:"metadata,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewDashboardMetrics creates a new DashboardMetrics instance
func NewDashboardMetrics() *DashboardMetrics {
	return &DashboardMetrics{}
}
