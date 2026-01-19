package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DeleteReviewByAdminInput represents input for DeleteReviewByAdmin
type DeleteReviewByAdminInput struct {
	ReviewId uuid.UUID
	AdminId uuid.UUID
	Reason string
}

// DeleteReviewByAdminUsecase defines the interface for DeleteReviewByAdmin
type DeleteReviewByAdminUsecase interface {
	Execute(ctx context.Context, input DeleteReviewByAdminInput) error
}

type delete_review_by_adminUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewDeleteReviewByAdminUsecase creates a new instance
func NewDeleteReviewByAdminUsecase() DeleteReviewByAdminUsecase {
	return &delete_review_by_adminUsecaseImpl{}
}

// Execute executes DeleteReviewByAdmin
func (u *delete_review_by_adminUsecaseImpl) Execute(ctx context.Context, input DeleteReviewByAdminInput) error {
	// TODO: Implement business logic

	return nil
}
