package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// PostReviewInput represents input for PostReview
type PostReviewInput struct {
	ProductId uuid.UUID
	OrderId uuid.UUID
	CustomerId uuid.UUID
	Nickname string
	Rating int
	Title string
	Content string
	ImageUrls []string
}

// PostReviewUsecase defines the interface for PostReview
type PostReviewUsecase interface {
	Execute(ctx context.Context, input PostReviewInput) error
}

type post_reviewUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewPostReviewUsecase creates a new instance
func NewPostReviewUsecase() PostReviewUsecase {
	return &post_reviewUsecaseImpl{}
}

// Execute executes PostReview
func (u *post_reviewUsecaseImpl) Execute(ctx context.Context, input PostReviewInput) error {
	// TODO: Implement business logic

	return nil
}
