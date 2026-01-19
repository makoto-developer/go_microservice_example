package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpdatePopularKeywordsInput represents input for UpdatePopularKeywords
type UpdatePopularKeywordsInput struct {
	PeriodType PeriodType
}

// UpdatePopularKeywordsUsecase defines the interface for UpdatePopularKeywords
type UpdatePopularKeywordsUsecase interface {
	Execute(ctx context.Context, input UpdatePopularKeywordsInput) error
}

type update_popular_keywordsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUpdatePopularKeywordsUsecase creates a new instance
func NewUpdatePopularKeywordsUsecase() UpdatePopularKeywordsUsecase {
	return &update_popular_keywordsUsecaseImpl{}
}

// Execute executes UpdatePopularKeywords
func (u *update_popular_keywordsUsecaseImpl) Execute(ctx context.Context, input UpdatePopularKeywordsInput) error {
	// TODO: Implement business logic

	return nil
}
