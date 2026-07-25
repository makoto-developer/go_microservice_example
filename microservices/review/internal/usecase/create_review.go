package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/review/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/review/internal/repository"
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
	Rating   int
}

type CreateReviewUsecase interface {
	Execute(ctx context.Context, input CreateReviewInput) (CreateReviewOutput, error)
}

type createReviewUsecaseImpl struct {
	reviewRepo repository.ReviewRepository
}

func NewCreateReviewUsecase(reviewRepo repository.ReviewRepository) CreateReviewUsecase {
	return &createReviewUsecaseImpl{
		reviewRepo: reviewRepo,
	}
}

func (u *createReviewUsecaseImpl) Execute(ctx context.Context, input CreateReviewInput) (CreateReviewOutput, error) {
	review, err := domain.NewReview(
		input.CustomerID,
		input.ProductID,
		input.OrderID,
		input.Rating,
		input.ReviewText,
	)
	if err != nil {
		return CreateReviewOutput{}, err
	}

	err = u.reviewRepo.Create(ctx, review)
	if err != nil {
		return CreateReviewOutput{}, err
	}

	return CreateReviewOutput{
		ReviewID: review.ID,
		Rating:   review.Rating,
	}, nil
}
