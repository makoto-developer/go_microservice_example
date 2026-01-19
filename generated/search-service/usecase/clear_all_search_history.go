package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ClearAllSearchHistoryInput represents input for ClearAllSearchHistory
type ClearAllSearchHistoryInput struct {
	UserId uuid.UUID
}

// ClearAllSearchHistoryUsecase defines the interface for ClearAllSearchHistory
type ClearAllSearchHistoryUsecase interface {
	Execute(ctx context.Context, input ClearAllSearchHistoryInput) error
}

type clear_all_search_historyUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewClearAllSearchHistoryUsecase creates a new instance
func NewClearAllSearchHistoryUsecase() ClearAllSearchHistoryUsecase {
	return &clear_all_search_historyUsecaseImpl{}
}

// Execute executes ClearAllSearchHistory
func (u *clear_all_search_historyUsecaseImpl) Execute(ctx context.Context, input ClearAllSearchHistoryInput) error {
	// TODO: Implement business logic

	return nil
}
