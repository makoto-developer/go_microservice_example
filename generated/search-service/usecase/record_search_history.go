package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// RecordSearchHistoryInput represents input for RecordSearchHistory
type RecordSearchHistoryInput struct {
	UserId uuid.UUID
	Keyword string
	ResultCount int
	ClickedProductId uuid.UUID
}

// RecordSearchHistoryUsecase defines the interface for RecordSearchHistory
type RecordSearchHistoryUsecase interface {
	Execute(ctx context.Context, input RecordSearchHistoryInput) error
}

type record_search_historyUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewRecordSearchHistoryUsecase creates a new instance
func NewRecordSearchHistoryUsecase() RecordSearchHistoryUsecase {
	return &record_search_historyUsecaseImpl{}
}

// Execute executes RecordSearchHistory
func (u *record_search_historyUsecaseImpl) Execute(ctx context.Context, input RecordSearchHistoryInput) error {
	// TODO: Implement business logic

	return nil
}
