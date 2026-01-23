package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/search/internal/domain"
)

type SearchRepository interface {
	IndexEntity(ctx context.Context, index *domain.SearchIndex) error
	Search(ctx context.Context, query string, entityType domain.IndexType) ([]*domain.SearchResult, error)
	GetByEntityID(ctx context.Context, entityType domain.IndexType, entityID uuid.UUID) (*domain.SearchIndex, error)
	Update(ctx context.Context, index *domain.SearchIndex) error
	Delete(ctx context.Context, entityType domain.IndexType, entityID uuid.UUID) error
}
