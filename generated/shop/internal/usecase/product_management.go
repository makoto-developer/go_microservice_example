package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

// ProductCreateInput represents the input for product creation
type ProductCreateInput struct {
	ShopID      uuid.UUID
	Name        string
	Description string
	Price       int64
	CategoryID  uuid.UUID
	Tags        []string
	Weight      *float64
	Dimensions  *string
	JANCode     *string
	StockCount  int
}

// ProductUpdateInput represents the input for product update
type ProductUpdateInput struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       int64
	CategoryID  uuid.UUID
	Tags        []string
	Weight      *float64
	Dimensions  *string
	JANCode     *string
	StockCount  int
}

// ProductManagementUsecase handles product management operations
type ProductManagementUsecase interface {
	CreateProduct(ctx context.Context, input ProductCreateInput) (uuid.UUID, error)
	UpdateProduct(ctx context.Context, input ProductUpdateInput) error
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
	GetProduct(ctx context.Context, productID uuid.UUID) (*domain.Product, error)
	ListProductsByShop(ctx context.Context, shopID uuid.UUID) ([]*domain.Product, error)
	PublishProduct(ctx context.Context, productID uuid.UUID) error
	UnpublishProduct(ctx context.Context, productID uuid.UUID) error
}

type productManagementUsecase struct {
	productRepo repository.ProductRepository
	shopRepo    repository.ShopRepository
}

// NewProductManagementUsecase creates a new product management usecase
func NewProductManagementUsecase(
	productRepo repository.ProductRepository,
	shopRepo repository.ShopRepository,
) ProductManagementUsecase {
	return &productManagementUsecase{
		productRepo: productRepo,
		shopRepo:    shopRepo,
	}
}

func (u *productManagementUsecase) CreateProduct(ctx context.Context, input ProductCreateInput) (uuid.UUID, error) {
	// Validate shop exists
	shop, err := u.shopRepo.GetByID(ctx, input.ShopID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("shop not found: %w", err)
	}

	// Check if shop is approved
	if shop.Status != domain.ShopStatusApproved {
		return uuid.Nil, fmt.Errorf("shop is not approved")
	}

	// Validate input
	if err := u.validateCreateInput(input); err != nil {
		return uuid.Nil, err
	}

	// Create product
	product := &domain.Product{
		ID:          uuid.New(),
		ShopID:      input.ShopID,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
		Tags:        input.Tags,
		Weight:      input.Weight,
		Dimensions:  input.Dimensions,
		JANCode:     input.JANCode,
		StockCount:  input.StockCount,
		Status:      domain.ProductStatusDraft,
		IsPublic:    false,
		IsDeleted:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product.ID, nil
}

func (u *productManagementUsecase) UpdateProduct(ctx context.Context, input ProductUpdateInput) error {
	// Validate input
	if err := u.validateUpdateInput(input); err != nil {
		return err
	}

	// Get existing product
	product, err := u.productRepo.GetByID(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Update product fields
	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.CategoryID = input.CategoryID
	product.Tags = input.Tags
	product.Weight = input.Weight
	product.Dimensions = input.Dimensions
	product.JANCode = input.JANCode
	product.StockCount = input.StockCount
	product.UpdatedAt = time.Now()

	if err := u.productRepo.Update(ctx, product); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (u *productManagementUsecase) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	if err := u.productRepo.Delete(ctx, productID); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

func (u *productManagementUsecase) GetProduct(ctx context.Context, productID uuid.UUID) (*domain.Product, error) {
	product, err := u.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (u *productManagementUsecase) ListProductsByShop(ctx context.Context, shopID uuid.UUID) ([]*domain.Product, error) {
	products, err := u.productRepo.GetByShopID(ctx, shopID, false)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}

func (u *productManagementUsecase) PublishProduct(ctx context.Context, productID uuid.UUID) error {
	// Get product to validate
	_, err := u.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Update status
	if err := u.productRepo.UpdateStatus(ctx, productID, domain.ProductStatusPublished); err != nil {
		return fmt.Errorf("failed to publish product: %w", err)
	}

	// Set public
	if err := u.productRepo.UpdateIsPublic(ctx, productID, true); err != nil {
		return fmt.Errorf("failed to set product public: %w", err)
	}

	return nil
}

func (u *productManagementUsecase) UnpublishProduct(ctx context.Context, productID uuid.UUID) error {
	if err := u.productRepo.UpdateIsPublic(ctx, productID, false); err != nil {
		return fmt.Errorf("failed to unpublish product: %w", err)
	}
	return nil
}

func (u *productManagementUsecase) validateCreateInput(input ProductCreateInput) error {
	if input.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if input.Price <= 0 {
		return fmt.Errorf("product price must be greater than 0")
	}
	if input.StockCount < 0 {
		return fmt.Errorf("stock count cannot be negative")
	}
	return nil
}

func (u *productManagementUsecase) validateUpdateInput(input ProductUpdateInput) error {
	if input.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if input.Price <= 0 {
		return fmt.Errorf("product price must be greater than 0")
	}
	if input.StockCount < 0 {
		return fmt.Errorf("stock count cannot be negative")
	}
	return nil
}
