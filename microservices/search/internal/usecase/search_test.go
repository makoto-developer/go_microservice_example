package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/search/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/search/internal/usecase"
)

type mockSearchRepository struct {
	searchFunc func(ctx context.Context, query string, entityType domain.IndexType) ([]*domain.SearchResult, error)
}

func (m *mockSearchRepository) IndexEntity(ctx context.Context, index *domain.SearchIndex) error {
	return nil
}

func (m *mockSearchRepository) Search(ctx context.Context, query string, entityType domain.IndexType) ([]*domain.SearchResult, error) {
	if m.searchFunc != nil {
		return m.searchFunc(ctx, query, entityType)
	}
	return nil, nil
}

func (m *mockSearchRepository) GetByEntityID(ctx context.Context, entityType domain.IndexType, entityID uuid.UUID) (*domain.SearchIndex, error) {
	return nil, nil
}

func (m *mockSearchRepository) Update(ctx context.Context, index *domain.SearchIndex) error {
	return nil
}

func (m *mockSearchRepository) Delete(ctx context.Context, entityType domain.IndexType, entityID uuid.UUID) error {
	return nil
}

func TestSearchUsecase_Success(t *testing.T) {
	expectedResults := []*domain.SearchResult{
		{
			EntityType:  domain.IndexTypeProduct,
			EntityID:    uuid.New(),
			Title:       "Test Product",
			Description: "Test Description",
			Score:       0.95,
		},
	}

	repo := &mockSearchRepository{
		searchFunc: func(ctx context.Context, query string, entityType domain.IndexType) ([]*domain.SearchResult, error) {
			if query != "test" {
				t.Errorf("expected query 'test', got %v", query)
			}
			return expectedResults, nil
		},
	}

	uc := usecase.NewSearchUsecase(repo)

	input := usecase.SearchInput{
		Query:      "test",
		EntityType: domain.IndexTypeProduct,
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.Total != 1 {
		t.Errorf("expected total 1, got %d", output.Total)
	}

	if len(output.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(output.Results))
	}
}

func TestSearchUsecase_NoResults(t *testing.T) {
	repo := &mockSearchRepository{
		searchFunc: func(ctx context.Context, query string, entityType domain.IndexType) ([]*domain.SearchResult, error) {
			return []*domain.SearchResult{}, nil
		},
	}

	uc := usecase.NewSearchUsecase(repo)

	input := usecase.SearchInput{
		Query:      "nonexistent",
		EntityType: domain.IndexTypeProduct,
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.Total != 0 {
		t.Errorf("expected total 0, got %d", output.Total)
	}
}
