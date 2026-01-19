package domain

import (
	"github.com/google/uuid"
	"time"
)

// SagaState represents SagaState
type SagaState struct {
	Id uuid.UUID `db:"id" json:"id"`
	OrderId uuid.UUID `db:"order_id" json:"order_id"`
	SagaType SagaType `db:"saga_type" json:"saga_type"`
	CurrentStep string `db:"current_step" json:"current_step"`
	Status SagaStatus `db:"status" json:"status"`
	CompensationData *map[string]interface{} `db:"compensation_data" json:"compensation_data,omitempty"`
	ErrorMessage *text `db:"error_message" json:"error_message,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewSagaState creates a new SagaState instance
func NewSagaState() *SagaState {
	return &SagaState{}
}
