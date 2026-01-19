package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/usecase"
)

type mockCustomerRepository struct {
	getByIDFunc func(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
}

func (m *mockCustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	return nil
}

func (m *mockCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	return nil
}

func (m *mockCustomerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Customer, error) {
	return nil, nil
}

func TestGetCustomerProfileUsecase_Success(t *testing.T) {
	customerID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	expectedCustomer := &domain.Customer{
		ID:        customerID,
		UserID:    userID,
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "090-1234-5678",
		CreatedAt: now,
		UpdatedAt: now,
	}

	repo := &mockCustomerRepository{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
			if id != customerID {
				t.Errorf("expected customer ID %v, got %v", customerID, id)
			}
			return expectedCustomer, nil
		},
	}

	uc := usecase.NewGetCustomerProfileUsecase(repo)

	input := usecase.GetCustomerProfileInput{CustomerID: customerID}
	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.Customer.ID != expectedCustomer.ID {
		t.Errorf("expected ID %v, got %v", expectedCustomer.ID, output.Customer.ID)
	}

	if output.Customer.FirstName != expectedCustomer.FirstName {
		t.Errorf("expected FirstName %v, got %v", expectedCustomer.FirstName, output.Customer.FirstName)
	}
}

func TestGetCustomerProfileUsecase_NotFound(t *testing.T) {
	customerID := uuid.New()

	repo := &mockCustomerRepository{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
			return nil, domain.ErrCustomerNotFound
		},
	}

	uc := usecase.NewGetCustomerProfileUsecase(repo)

	input := usecase.GetCustomerProfileInput{CustomerID: customerID}
	_, err := uc.Execute(context.Background(), input)

	if err != domain.ErrCustomerNotFound {
		t.Errorf("expected ErrCustomerNotFound, got %v", err)
	}
}

func TestGetCustomerProfileUsecase_RepositoryError(t *testing.T) {
	customerID := uuid.New()
	expectedErr := errors.New("database error")

	repo := &mockCustomerRepository{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
			return nil, expectedErr
		},
	}

	uc := usecase.NewGetCustomerProfileUsecase(repo)

	input := usecase.GetCustomerProfileInput{CustomerID: customerID}
	_, err := uc.Execute(context.Background(), input)

	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}
