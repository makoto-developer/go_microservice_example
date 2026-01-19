package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/usecase"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository for testing
type MockProductRepositoryForGetProduct struct {
	mock.Mock
}

func (m *MockProductRepositoryForGetProduct) Create(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepositoryForGetProduct) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepositoryForGetProduct) Update(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepositoryForGetProduct) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepositoryForGetProduct) List(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *MockProductRepositoryForGetProduct) ListByShopID(ctx context.Context, shopID uuid.UUID, category string, publishedOnly bool, limit, offset int) ([]*domain.Product, error) {
	args := m.Called(ctx, shopID, category, publishedOnly, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func TestGetProductUsecase_Execute_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProductRepositoryForGetProduct)
	uc := usecase.NewGetProductUsecase(mockRepo)

	productID := uuid.New()
	shopID := uuid.New()
	price := decimal.NewFromInt(29800)
	weight := decimal.NewFromFloat(0.055)
	size := "5.4 x 4.6 x 2.1 cm"
	janCode := "4901234567890"

	expectedProduct := &domain.Product{
		Id:            productID,
		ShopId:        shopID,
		Name:          "ワイヤレスイヤホン Pro",
		Description:   "ノイズキャンセリング機能付きの高音質ワイヤレスイヤホン",
		Price:         price,
		Category:      "オーディオ",
		StockQuantity: 50,
		Weight:        &weight,
		Size:          &size,
		JanCode:       &janCode,
		Published:     true,
		Deleted:       false,
	}

	mockRepo.On("FindByID", ctx, productID).Return(expectedProduct, nil)

	// Act
	output, err := uc.Execute(ctx, usecase.GetProductInput{
		ProductID: productID,
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, expectedProduct, output.Product)
	mockRepo.AssertExpectations(t)
}

func TestGetProductUsecase_Execute_ProductNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProductRepositoryForGetProduct)
	uc := usecase.NewGetProductUsecase(mockRepo)

	productID := uuid.New()

	mockRepo.On("FindByID", ctx, productID).Return(nil, fmt.Errorf("product not found: %s", productID))

	// Act
	output, err := uc.Execute(ctx, usecase.GetProductInput{
		ProductID: productID,
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "failed to find product")
	mockRepo.AssertExpectations(t)
}

func TestGetProductUsecase_Execute_ProductDeleted(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProductRepositoryForGetProduct)
	uc := usecase.NewGetProductUsecase(mockRepo)

	productID := uuid.New()
	shopID := uuid.New()
	price := decimal.NewFromInt(29800)

	deletedProduct := &domain.Product{
		Id:            productID,
		ShopId:        shopID,
		Name:          "削除された商品",
		Description:   "この商品は削除されました",
		Price:         price,
		Category:      "その他",
		StockQuantity: 0,
		Published:     false,
		Deleted:       true, // 削除済み
	}

	mockRepo.On("FindByID", ctx, productID).Return(deletedProduct, nil)

	// Act
	output, err := uc.Execute(ctx, usecase.GetProductInput{
		ProductID: productID,
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "product not found (deleted)")
	mockRepo.AssertExpectations(t)
}

func TestGetProductUsecase_Execute_EmptyProductID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProductRepositoryForGetProduct)
	uc := usecase.NewGetProductUsecase(mockRepo)

	// Act
	output, err := uc.Execute(ctx, usecase.GetProductInput{
		ProductID: uuid.Nil, // 空のUUID
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "product_id is required")
	mockRepo.AssertExpectations(t)
}

func TestGetProductUsecase_Execute_RepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProductRepositoryForGetProduct)
	uc := usecase.NewGetProductUsecase(mockRepo)

	productID := uuid.New()

	mockRepo.On("FindByID", ctx, productID).Return(nil, fmt.Errorf("database connection error"))

	// Act
	output, err := uc.Execute(ctx, usecase.GetProductInput{
		ProductID: productID,
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "failed to find product")
	mockRepo.AssertExpectations(t)
}
