package usecase

import (
	"context"

	"github.com/makoto-developer/go_microservice_example/microservices/search/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/search/internal/repository"
)

type SearchInput struct {
	Query      string
	EntityType domain.IndexType
}

type SearchOutput struct {
	Results []*domain.SearchResult
	Total   int
}

type SearchUsecase interface {
	Execute(ctx context.Context, input SearchInput) (SearchOutput, error)
}

type searchUsecaseImpl struct {
	searchRepo repository.SearchRepository
}

func NewSearchUsecase(searchRepo repository.SearchRepository) SearchUsecase {
	return &searchUsecaseImpl{
		searchRepo: searchRepo,
	}
}

func (u *searchUsecaseImpl) Execute(ctx context.Context, input SearchInput) (SearchOutput, error) {
	results, err := u.searchRepo.Search(ctx, input.Query, input.EntityType)
	if err != nil {
		return SearchOutput{}, err
	}

	return SearchOutput{
		Results: results,
		Total:   len(results),
	}, nil
}
