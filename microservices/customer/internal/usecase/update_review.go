package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type UpdateReviewInput struct {
	CustomerID uuid.UUID
	ReviewID   uuid.UUID
	Rating     int
	ReviewText string
}

type UpdateReviewOutput struct {
	Review *domain.Review
}

type UpdateReviewUsecase interface {
	Execute(ctx context.Context, input UpdateReviewInput) (UpdateReviewOutput, error)
}

type updateReviewUsecase struct {
	reviewRepo repository.ReviewRepository
}

func NewUpdateReviewUsecase(reviewRepo repository.ReviewRepository) UpdateReviewUsecase {
	return &updateReviewUsecase{reviewRepo: reviewRepo}
}

func (u *updateReviewUsecase) Execute(ctx context.Context, input UpdateReviewInput) (UpdateReviewOutput, error) {
	review, err := u.reviewRepo.GetByID(ctx, input.ReviewID)
	if err != nil {
		return UpdateReviewOutput{}, err
	}

	if review.CustomerID != input.CustomerID {
		return UpdateReviewOutput{}, domain.ErrReviewNotFound
	}

	if input.Rating < 1 || input.Rating > 5 {
		return UpdateReviewOutput{}, domain.ErrInvalidRating
	}

	review.Rating = input.Rating
	review.ReviewText = input.ReviewText
	review.UpdatedAt = time.Now()

	if err := u.reviewRepo.Update(ctx, review); err != nil {
		return UpdateReviewOutput{}, err
	}

	return UpdateReviewOutput{Review: review}, nil
}
