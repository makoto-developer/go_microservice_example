package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UnmarkReviewHelpfulInput represents input for UnmarkReviewHelpful
type UnmarkReviewHelpfulInput struct {
	ReviewId uuid.UUID
	UserId uuid.UUID
}

// UnmarkReviewHelpfulUsecase defines the interface for UnmarkReviewHelpful
type UnmarkReviewHelpfulUsecase interface {
	Execute(ctx context.Context, input UnmarkReviewHelpfulInput) error
}

type unmark_review_helpfulUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUnmarkReviewHelpfulUsecase creates a new instance
func NewUnmarkReviewHelpfulUsecase() UnmarkReviewHelpfulUsecase {
	return &unmark_review_helpfulUsecaseImpl{}
}

// Execute executes UnmarkReviewHelpful
func (u *unmark_review_helpfulUsecaseImpl) Execute(ctx context.Context, input UnmarkReviewHelpfulInput) error {
	// TODO: Implement business logic

	return nil
}
