package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
	"github.com/shopspring/decimal"
)

type RegisterProductInput struct {
	ShopID        uuid.UUID
	Name          string
	Description   string
	Price         decimal.Decimal
	Category      string
	StockQuantity int
	Weight        *decimal.Decimal
	Size          *string
	JanCode       *string
	Tags          []string
}

type RegisterProductOutput struct {
	ProductID uuid.UUID
	Message   string
}

type RegisterProductUsecase interface {
	Execute(ctx context.Context, input RegisterProductInput) (*RegisterProductOutput, error)
}

type registerProductUsecaseImpl struct {
	productRepo    domain.ProductRepository
	productTagRepo domain.ProductTagRepository
}

func NewRegisterProductUsecase(
	productRepo domain.ProductRepository,
	productTagRepo domain.ProductTagRepository,
) RegisterProductUsecase {
	return &registerProductUsecaseImpl{
		productRepo:    productRepo,
		productTagRepo: productTagRepo,
	}
}

func (u *registerProductUsecaseImpl) Execute(ctx context.Context, input RegisterProductInput) (*RegisterProductOutput, error) {
	product := &domain.Product{
		Id:            uuid.New(),
		ShopId:        input.ShopID,
		Name:          input.Name,
		Description:   input.Description,
		Price:         input.Price,
		Category:      input.Category,
		StockQuantity: input.StockQuantity,
		Weight:        input.Weight,
		Size:          input.Size,
		JanCode:       input.JanCode,
		Published:     false,
		Deleted:       false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	for _, tagName := range input.Tags {
		tag := &domain.ProductTag{
			Id:        uuid.New(),
			ProductId: product.Id,
			TagName:   tagName,
			CreatedAt: time.Now(),
		}
		if err := u.productTagRepo.Create(ctx, tag); err != nil {
			return nil, fmt.Errorf("failed to create product tag: %w", err)
		}
	}

	return &RegisterProductOutput{
		ProductID: product.Id,
		Message:   "Product registered successfully",
	}, nil
}
