package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetSearchHistoryInput represents input for GetSearchHistory
type GetSearchHistoryInput struct {
	UserId uuid.UUID
	Limit int
}

// GetSearchHistoryUsecase defines the interface for GetSearchHistory
type GetSearchHistoryUsecase interface {
	Execute(ctx context.Context, input GetSearchHistoryInput) error
}

type get_search_historyUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetSearchHistoryUsecase creates a new instance
func NewGetSearchHistoryUsecase() GetSearchHistoryUsecase {
	return &get_search_historyUsecaseImpl{}
}

// Execute executes GetSearchHistory
func (u *get_search_historyUsecaseImpl) Execute(ctx context.Context, input GetSearchHistoryInput) error {
	// TODO: Implement business logic

	return nil
}
