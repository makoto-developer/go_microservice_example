package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository"
)

type CreateReviewInput struct {
	CustomerID uuid.UUID
	ProductID  uuid.UUID
	OrderID    uuid.UUID
	Rating     int
	ReviewText string
}

type CreateReviewOutput struct {
	ReviewID uuid.UUID
}

type CreateReviewUsecase interface {
	Execute(ctx context.Context, input CreateReviewInput) (CreateReviewOutput, error)
}

type createReviewUsecase struct {
	reviewRepo repository.ReviewRepository
}

func NewCreateReviewUsecase(reviewRepo repository.ReviewRepository) CreateReviewUsecase {
	return &createReviewUsecase{reviewRepo: reviewRepo}
}

func (u *createReviewUsecase) Execute(ctx context.Context, input CreateReviewInput) (CreateReviewOutput, error) {
	if input.Rating < 1 || input.Rating > 5 {
		return CreateReviewOutput{}, domain.ErrInvalidRating
	}

	review := domain.NewReview(
		input.CustomerID,
		input.ProductID,
		input.OrderID,
		input.Rating,
		input.ReviewText,
	)

	if err := u.reviewRepo.Create(ctx, review); err != nil {
		return CreateReviewOutput{}, err
	}

	return CreateReviewOutput{ReviewID: review.ID}, nil
}
