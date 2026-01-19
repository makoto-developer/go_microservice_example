package domain

import (
	"github.com/google/uuid"
	"time"
)

// PushTemplate represents PushTemplate
type PushTemplate struct {
	Id uuid.UUID `db:"id" json:"id"`
	TemplateKey string `db:"template_key" json:"template_key"`
	TitleTemplate string `db:"title_template" json:"title_template"`
	BodyTemplate string `db:"body_template" json:"body_template"`
	Variables []string `db:"variables" json:"variables"`
	Version int `db:"version" json:"version"`
	IsActive bool `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewPushTemplate creates a new PushTemplate instance
func NewPushTemplate() *PushTemplate {
	return &PushTemplate{}
}
