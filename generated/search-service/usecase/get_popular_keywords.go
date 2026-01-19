package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetPopularKeywordsInput represents input for GetPopularKeywords
type GetPopularKeywordsInput struct {
	PeriodType PeriodType
	Category string
	Limit int
}

// GetPopularKeywordsUsecase defines the interface for GetPopularKeywords
type GetPopularKeywordsUsecase interface {
	Execute(ctx context.Context, input GetPopularKeywordsInput) error
}

type get_popular_keywordsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetPopularKeywordsUsecase creates a new instance
func NewGetPopularKeywordsUsecase() GetPopularKeywordsUsecase {
	return &get_popular_keywordsUsecaseImpl{}
}

// Execute executes GetPopularKeywords
func (u *get_popular_keywordsUsecaseImpl) Execute(ctx context.Context, input GetPopularKeywordsInput) error {
	// TODO: Implement business logic

	return nil
}
