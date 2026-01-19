package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ApproveReviewInput represents input for ApproveReview
type ApproveReviewInput struct {
	ReviewId uuid.UUID
	AdminId uuid.UUID
}

// ApproveReviewUsecase defines the interface for ApproveReview
type ApproveReviewUsecase interface {
	Execute(ctx context.Context, input ApproveReviewInput) error
}

type approve_reviewUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewApproveReviewUsecase creates a new instance
func NewApproveReviewUsecase() ApproveReviewUsecase {
	return &approve_reviewUsecaseImpl{}
}

// Execute executes ApproveReview
func (u *approve_reviewUsecaseImpl) Execute(ctx context.Context, input ApproveReviewInput) error {
	// TODO: Implement business logic

	return nil
}
