package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *domain.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error)
	GetByOrderAndProduct(ctx context.Context, orderID, productID uuid.UUID) (*domain.Review, error)
	Update(ctx context.Context, review *domain.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.Review, error)
	ListByCustomer(ctx context.Context, customerID uuid.UUID) ([]*domain.Review, error)

	AddImage(ctx context.Context, image *domain.ReviewImage) error
	GetImages(ctx context.Context, reviewID uuid.UUID) ([]*domain.ReviewImage, error)
}
