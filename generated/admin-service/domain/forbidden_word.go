package domain

import (
	"time"
	"github.com/google/uuid"
)

// ForbiddenWord represents ForbiddenWord
type ForbiddenWord struct {
	Id uuid.UUID `db:"id" json:"id"`
	Word string `db:"word" json:"word"`
	Context WordContext `db:"context" json:"context"`
	Severity Severity `db:"severity" json:"severity"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewForbiddenWord creates a new ForbiddenWord instance
func NewForbiddenWord() *ForbiddenWord {
	return &ForbiddenWord{}
}
