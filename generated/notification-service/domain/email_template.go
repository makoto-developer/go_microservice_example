package domain

import (
	"github.com/google/uuid"
	"time"
)

// EmailTemplate represents EmailTemplate
type EmailTemplate struct {
	Id uuid.UUID `db:"id" json:"id"`
	TemplateKey string `db:"template_key" json:"template_key"`
	SubjectTemplate string `db:"subject_template" json:"subject_template"`
	HtmlTemplate text `db:"html_template" json:"html_template"`
	TextTemplate text `db:"text_template" json:"text_template"`
	Variables []string `db:"variables" json:"variables"`
	Version int `db:"version" json:"version"`
	IsActive bool `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewEmailTemplate creates a new EmailTemplate instance
func NewEmailTemplate() *EmailTemplate {
	return &EmailTemplate{}
}
