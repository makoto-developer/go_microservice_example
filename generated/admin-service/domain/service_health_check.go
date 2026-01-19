package domain

import (
	"github.com/google/uuid"
	"time"
)

// ServiceHealthCheck represents ServiceHealthCheck
type ServiceHealthCheck struct {
	Id uuid.UUID `db:"id" json:"id"`
	ServiceName string `db:"service_name" json:"service_name"`
	Status HealthStatus `db:"status" json:"status"`
	ResponseTimeMs int `db:"response_time_ms" json:"response_time_ms"`
	ErrorMessage *text `db:"error_message" json:"error_message,omitempty"`
	CheckedAt time.Time `db:"checked_at" json:"checked_at"`
}

// NewServiceHealthCheck creates a new ServiceHealthCheck instance
func NewServiceHealthCheck() *ServiceHealthCheck {
	return &ServiceHealthCheck{}
}
