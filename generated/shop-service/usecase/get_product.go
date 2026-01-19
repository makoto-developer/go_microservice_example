package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type GetProductInput struct {
	ProductID uuid.UUID
}

type GetProductOutput struct {
	Product *domain.Product
}

type GetProductUsecase interface {
	Execute(ctx context.Context, input GetProductInput) (*GetProductOutput, error)
}

type getProductUsecaseImpl struct {
	productRepo domain.ProductRepository
}

func NewGetProductUsecase(productRepo domain.ProductRepository) GetProductUsecase {
	return &getProductUsecaseImpl{
		productRepo: productRepo,
	}
}

func (u *getProductUsecaseImpl) Execute(ctx context.Context, input GetProductInput) (*GetProductOutput, error) {
	log.Printf("GetProductUsecase: product_id=%s", input.ProductID)

	// Validate input
	if input.ProductID == uuid.Nil {
		return nil, fmt.Errorf("product_id is required")
	}

	// Find product by ID
	product, err := u.productRepo.FindByID(ctx, input.ProductID)
	if err != nil {
		log.Printf("GetProductUsecase error: %v", err)
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	// Check if product is deleted
	if product.Deleted {
		return nil, fmt.Errorf("product not found (deleted): %s", input.ProductID)
	}

	log.Printf("GetProductUsecase: found product %s", product.Name)
	return &GetProductOutput{
		Product: product,
	}, nil
}
