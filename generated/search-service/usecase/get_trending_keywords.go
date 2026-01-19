package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetTrendingKeywordsInput represents input for GetTrendingKeywords
type GetTrendingKeywordsInput struct {
	Limit int
}

// GetTrendingKeywordsUsecase defines the interface for GetTrendingKeywords
type GetTrendingKeywordsUsecase interface {
	Execute(ctx context.Context, input GetTrendingKeywordsInput) error
}

type get_trending_keywordsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetTrendingKeywordsUsecase creates a new instance
func NewGetTrendingKeywordsUsecase() GetTrendingKeywordsUsecase {
	return &get_trending_keywordsUsecaseImpl{}
}

// Execute executes GetTrendingKeywords
func (u *get_trending_keywordsUsecaseImpl) Execute(ctx context.Context, input GetTrendingKeywordsInput) error {
	// TODO: Implement business logic

	return nil
}
