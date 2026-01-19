package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DeleteSearchHistoryInput represents input for DeleteSearchHistory
type DeleteSearchHistoryInput struct {
	UserId uuid.UUID
	HistoryId uuid.UUID
}

// DeleteSearchHistoryUsecase defines the interface for DeleteSearchHistory
type DeleteSearchHistoryUsecase interface {
	Execute(ctx context.Context, input DeleteSearchHistoryInput) error
}

type delete_search_historyUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewDeleteSearchHistoryUsecase creates a new instance
func NewDeleteSearchHistoryUsecase() DeleteSearchHistoryUsecase {
	return &delete_search_historyUsecaseImpl{}
}

// Execute executes DeleteSearchHistory
func (u *delete_search_historyUsecaseImpl) Execute(ctx context.Context, input DeleteSearchHistoryInput) error {
	// TODO: Implement business logic

	return nil
}
