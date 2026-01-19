package domain

import (
	"github.com/google/uuid"
	"time"
)

// Category represents Category
type Category struct {
	Id uuid.UUID `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	ParentId *uuid.UUID `db:"parent_id" json:"parent_id,omitempty"`
	Level int `db:"level" json:"level"`
	DisplayOrder int `db:"display_order" json:"display_order"`
	IsActive bool `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewCategory creates a new Category instance
func NewCategory() *Category {
	return &Category{}
}
