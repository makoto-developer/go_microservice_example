package domain

import (
	"github.com/google/uuid"
	"time"
)

// AuditLog represents AuditLog
type AuditLog struct {
	Id uuid.UUID `db:"id" json:"id"`
	OperationType OperationType `db:"operation_type" json:"operation_type"`
	OperatorId uuid.UUID `db:"operator_id" json:"operator_id"`
	OperatorName string `db:"operator_name" json:"operator_name"`
	TargetType string `db:"target_type" json:"target_type"`
	TargetId uuid.UUID `db:"target_id" json:"target_id"`
	OperationDetail map[string]interface{} `db:"operation_detail" json:"operation_detail"`
	IpAddress string `db:"ip_address" json:"ip_address"`
	UserAgent *string `db:"user_agent" json:"user_agent,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewAuditLog creates a new AuditLog instance
func NewAuditLog() *AuditLog {
	return &AuditLog{}
}
