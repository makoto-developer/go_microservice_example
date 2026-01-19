package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type RegisterProductInput struct {
	ShopID        uuid.UUID
	Name          string
	Description   string
	Price         float64
	Category      string
	StockQuantity int
	Weight        *float64
	Size          *string
	JANCode       *string
	Tags          []string
}

type RegisterProductOutput struct {
	ProductID uuid.UUID
}

type RegisterProductUsecase interface {
	Execute(ctx context.Context, input RegisterProductInput) (RegisterProductOutput, error)
}

type registerProductUsecase struct {
	shopRepo    repository.ShopRepository
	productRepo repository.ProductRepository
}

func NewRegisterProductUsecase(shopRepo repository.ShopRepository, productRepo repository.ProductRepository) RegisterProductUsecase {
	return &registerProductUsecase{
		shopRepo:    shopRepo,
		productRepo: productRepo,
	}
}

func (u *registerProductUsecase) Execute(ctx context.Context, input RegisterProductInput) (RegisterProductOutput, error) {
	if _, err := u.shopRepo.GetByID(ctx, input.ShopID); err != nil {
		return RegisterProductOutput{}, err
	}

	if input.Name == "" || input.Price <= 0 {
		return RegisterProductOutput{}, domain.ErrInvalidProductData
	}

	product := domain.NewProduct(input.ShopID, input.Name, input.Description, input.Price, input.Category, input.StockQuantity)
	product.Weight = input.Weight
	product.Size = input.Size
	product.JANCode = input.JANCode

	if err := u.productRepo.Create(ctx, product); err != nil {
		return RegisterProductOutput{}, err
	}

	for _, tagName := range input.Tags {
		tag := domain.NewProductTag(product.ID, tagName)
		if err := u.productRepo.AddTag(ctx, tag); err != nil {
			return RegisterProductOutput{}, err
		}
	}

	return RegisterProductOutput{ProductID: product.ID}, nil
}
