package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetSearchSuggestionsInput represents input for GetSearchSuggestions
type GetSearchSuggestionsInput struct {
	Prefix string
	Limit int
}

// GetSearchSuggestionsUsecase defines the interface for GetSearchSuggestions
type GetSearchSuggestionsUsecase interface {
	Execute(ctx context.Context, input GetSearchSuggestionsInput) error
}

type get_search_suggestionsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetSearchSuggestionsUsecase creates a new instance
func NewGetSearchSuggestionsUsecase() GetSearchSuggestionsUsecase {
	return &get_search_suggestionsUsecaseImpl{}
}

// Execute executes GetSearchSuggestions
func (u *get_search_suggestionsUsecaseImpl) Execute(ctx context.Context, input GetSearchSuggestionsInput) error {
	// TODO: Implement business logic

	return nil
}
