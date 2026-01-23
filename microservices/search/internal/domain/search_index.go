package domain

import (
	"time"

	"github.com/google/uuid"
)

type IndexType string

const (
	IndexTypeProduct IndexType = "product"
	IndexTypeShop    IndexType = "shop"
)

type SearchIndex struct {
	ID          uuid.UUID `db:"id" json:"id"`
	EntityType  IndexType `db:"entity_type" json:"entity_type"`
	EntityID    uuid.UUID `db:"entity_id" json:"entity_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Keywords    string    `db:"keywords" json:"keywords"` // Space-separated keywords
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type SearchResult struct {
	EntityType  IndexType
	EntityID    uuid.UUID
	Title       string
	Description string
	Score       float64
}

func NewSearchIndex(entityType IndexType, entityID uuid.UUID, title, description, keywords string) *SearchIndex {
	now := time.Now()
	return &SearchIndex{
		ID:          uuid.New(),
		EntityType:  entityType,
		EntityID:    entityID,
		Title:       title,
		Description: description,
		Keywords:    keywords,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (s *SearchIndex) UpdateContent(title, description, keywords string) {
	s.Title = title
	s.Description = description
	s.Keywords = keywords
	s.UpdatedAt = time.Now()
}
