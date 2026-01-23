package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/review/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/review/internal/usecase"
)

type mockReviewRepository struct {
	createFunc func(ctx context.Context, review *domain.Review) error
}

func (m *mockReviewRepository) Create(ctx context.Context, review *domain.Review) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, review)
	}
	return nil
}

func (m *mockReviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error) {
	return nil, nil
}

func (m *mockReviewRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]*domain.Review, error) {
	return nil, nil
}

func (m *mockReviewRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Review, error) {
	return nil, nil
}

func (m *mockReviewRepository) Update(ctx context.Context, review *domain.Review) error {
	return nil
}

func (m *mockReviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func TestCreateReviewUsecase_Success(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()
	orderID := uuid.New()

	repo := &mockReviewRepository{
		createFunc: func(ctx context.Context, review *domain.Review) error {
			if review.CustomerID != customerID {
				t.Errorf("expected customer ID %v, got %v", customerID, review.CustomerID)
			}
			if review.Rating != 5 {
				t.Errorf("expected rating 5, got %d", review.Rating)
			}
			return nil
		},
	}

	uc := usecase.NewCreateReviewUsecase(repo)

	input := usecase.CreateReviewInput{
		CustomerID: customerID,
		ProductID:  productID,
		OrderID:    orderID,
		Rating:     5,
		ReviewText: "Great product!",
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.Rating != 5 {
		t.Errorf("expected rating 5, got %d", output.Rating)
	}
}

func TestCreateReviewUsecase_InvalidRating(t *testing.T) {
	repo := &mockReviewRepository{}
	uc := usecase.NewCreateReviewUsecase(repo)

	input := usecase.CreateReviewInput{
		CustomerID: uuid.New(),
		ProductID:  uuid.New(),
		OrderID:    uuid.New(),
		Rating:     6, // Invalid
		ReviewText: "Test",
	}

	_, err := uc.Execute(context.Background(), input)

	if err != domain.ErrInvalidRating {
		t.Errorf("expected ErrInvalidRating, got %v", err)
	}
}
