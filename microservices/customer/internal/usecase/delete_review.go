package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type DeleteReviewInput struct {
	CustomerID uuid.UUID
	ReviewID   uuid.UUID
}

type DeleteReviewOutput struct {
	Success bool
}

type DeleteReviewUsecase interface {
	Execute(ctx context.Context, input DeleteReviewInput) (DeleteReviewOutput, error)
}

type deleteReviewUsecase struct {
	reviewRepo repository.ReviewRepository
}

func NewDeleteReviewUsecase(reviewRepo repository.ReviewRepository) DeleteReviewUsecase {
	return &deleteReviewUsecase{reviewRepo: reviewRepo}
}

func (u *deleteReviewUsecase) Execute(ctx context.Context, input DeleteReviewInput) (DeleteReviewOutput, error) {
	review, err := u.reviewRepo.GetByID(ctx, input.ReviewID)
	if err != nil {
		return DeleteReviewOutput{}, err
	}

	if review.CustomerID != input.CustomerID {
		return DeleteReviewOutput{}, domain.ErrReviewNotFound
	}

	if err := u.reviewRepo.Delete(ctx, input.ReviewID); err != nil {
		return DeleteReviewOutput{}, err
	}

	return DeleteReviewOutput{Success: true}, nil
}
