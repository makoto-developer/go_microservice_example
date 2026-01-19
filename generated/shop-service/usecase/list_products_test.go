package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/usecase"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) List(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *MockProductRepository) ListByShopID(ctx context.Context, shopID uuid.UUID, category string, publishedOnly bool, limit, offset int) ([]*domain.Product, error) {
	args := m.Called(ctx, shopID, category, publishedOnly, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func TestListProductsUsecase_Execute_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	uc := usecase.NewListProductsUsecase(mockRepo)

	ctx := context.Background()
	shopID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	category := "電化製品"
	publishedOnly := true
	limit := 10
	offset := 0

	price, _ := decimal.NewFromString("29800")
	weight, _ := decimal.NewFromString("0.055")
	size := "5.4 x 4.6 x 2.1 cm"
	janCode := "4901234567890"

	expectedProducts := []*domain.Product{
		{
			Id:            uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			ShopId:        shopID,
			Name:          "テスト商品1",
			Description:   "テスト商品1の説明",
			Price:         price,
			Category:      category,
			StockQuantity: 100,
			Weight:        &weight,
			Size:          &size,
			JanCode:       &janCode,
			Published:     true,
			Deleted:       false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	mockRepo.On("ListByShopID", ctx, shopID, category, publishedOnly, limit, offset).
		Return(expectedProducts, nil)

	input := usecase.ListProductsInput{
		ShopID:        shopID,
		Category:      category,
		PublishedOnly: publishedOnly,
		Limit:         limit,
		Offset:        offset,
	}

	// Act
	output, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, len(expectedProducts), len(output.Products))
	assert.Equal(t, expectedProducts[0].Id, output.Products[0].Id)
	assert.Equal(t, expectedProducts[0].Name, output.Products[0].Name)
	assert.Equal(t, expectedProducts[0].Price, output.Products[0].Price)
	assert.Equal(t, expectedProducts[0].Category, output.Products[0].Category)

	mockRepo.AssertExpectations(t)
}

func TestListProductsUsecase_Execute_EmptyShopID(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	uc := usecase.NewListProductsUsecase(mockRepo)

	ctx := context.Background()
	emptyShopID := uuid.Nil
	category := ""
	publishedOnly := false
	limit := 100
	offset := 0

	price1, _ := decimal.NewFromString("29800")
	price2, _ := decimal.NewFromString("45000")

	expectedProducts := []*domain.Product{
		{
			Id:            uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			ShopId:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Name:          "商品1",
			Description:   "商品1の説明",
			Price:         price1,
			Category:      "オーディオ",
			StockQuantity: 50,
			Published:     true,
			Deleted:       false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Id:            uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			ShopId:        uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			Name:          "商品2",
			Description:   "商品2の説明",
			Price:         price2,
			Category:      "ファッション",
			StockQuantity: 30,
			Published:     true,
			Deleted:       false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	mockRepo.On("ListByShopID", ctx, emptyShopID, category, publishedOnly, limit, offset).
		Return(expectedProducts, nil)

	input := usecase.ListProductsInput{
		ShopID:        emptyShopID,
		Category:      category,
		PublishedOnly: publishedOnly,
		Limit:         limit,
		Offset:        offset,
	}

	// Act
	output, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 2, len(output.Products))

	mockRepo.AssertExpectations(t)
}

func TestListProductsUsecase_Execute_WithCategoryFilter(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	uc := usecase.NewListProductsUsecase(mockRepo)

	ctx := context.Background()
	shopID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	category := "オーディオ"
	publishedOnly := true
	limit := 10
	offset := 0

	price, _ := decimal.NewFromString("29800")

	expectedProducts := []*domain.Product{
		{
			Id:            uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			ShopId:        shopID,
			Name:          "オーディオ商品",
			Price:         price,
			Category:      category,
			StockQuantity: 50,
			Published:     true,
			Deleted:       false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	mockRepo.On("ListByShopID", ctx, shopID, category, publishedOnly, limit, offset).
		Return(expectedProducts, nil)

	input := usecase.ListProductsInput{
		ShopID:        shopID,
		Category:      category,
		PublishedOnly: publishedOnly,
		Limit:         limit,
		Offset:        offset,
	}

	// Act
	output, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 1, len(output.Products))
	assert.Equal(t, category, output.Products[0].Category)

	mockRepo.AssertExpectations(t)
}

func TestListProductsUsecase_Execute_PublishedOnlyFilter(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	uc := usecase.NewListProductsUsecase(mockRepo)

	ctx := context.Background()
	shopID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	category := ""
	publishedOnly := true
	limit := 10
	offset := 0

	price, _ := decimal.NewFromString("29800")

	expectedProducts := []*domain.Product{
		{
			Id:            uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			ShopId:        shopID,
			Name:          "公開商品",
			Price:         price,
			Category:      "カテゴリ1",
			StockQuantity: 50,
			Published:     true,
			Deleted:       false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	mockRepo.On("ListByShopID", ctx, shopID, category, publishedOnly, limit, offset).
		Return(expectedProducts, nil)

	input := usecase.ListProductsInput{
		ShopID:        shopID,
		Category:      category,
		PublishedOnly: publishedOnly,
		Limit:         limit,
		Offset:        offset,
	}

	// Act
	output, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	for _, product := range output.Products {
		assert.True(t, product.Published)
	}

	mockRepo.AssertExpectations(t)
}

func TestListProductsUsecase_Execute_WithPagination(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	uc := usecase.NewListProductsUsecase(mockRepo)

	ctx := context.Background()
	shopID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	category := ""
	publishedOnly := false
	limit := 5
	offset := 10

	price, _ := decimal.NewFromString("29800")

	expectedProducts := []*domain.Product{
		{
			Id:            uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			ShopId:        shopID,
			Name:          "ページ2の商品",
			Price:         price,
			Category:      "カテゴリ",
			StockQuantity: 50,
			Published:     true,
			Deleted:       false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	mockRepo.On("ListByShopID", ctx, shopID, category, publishedOnly, limit, offset).
		Return(expectedProducts, nil)

	input := usecase.ListProductsInput{
		ShopID:        shopID,
		Category:      category,
		PublishedOnly: publishedOnly,
		Limit:         limit,
		Offset:        offset,
	}

	// Act
	output, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, len(expectedProducts), len(output.Products))

	mockRepo.AssertExpectations(t)
}

func TestListProductsUsecase_Execute_EmptyResult(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	uc := usecase.NewListProductsUsecase(mockRepo)

	ctx := context.Background()
	shopID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	category := "存在しないカテゴリ"
	publishedOnly := true
	limit := 10
	offset := 0

	expectedProducts := []*domain.Product{}

	mockRepo.On("ListByShopID", ctx, shopID, category, publishedOnly, limit, offset).
		Return(expectedProducts, nil)

	input := usecase.ListProductsInput{
		ShopID:        shopID,
		Category:      category,
		PublishedOnly: publishedOnly,
		Limit:         limit,
		Offset:        offset,
	}

	// Act
	output, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 0, len(output.Products))

	mockRepo.AssertExpectations(t)
}

func TestListProductsUsecase_Execute_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	uc := usecase.NewListProductsUsecase(mockRepo)

	ctx := context.Background()
	shopID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	category := ""
	publishedOnly := false
	limit := 10
	offset := 0

	expectedError := assert.AnError

	mockRepo.On("ListByShopID", ctx, shopID, category, publishedOnly, limit, offset).
		Return(nil, expectedError)

	input := usecase.ListProductsInput{
		ShopID:        shopID,
		Category:      category,
		PublishedOnly: publishedOnly,
		Limit:         limit,
		Offset:        offset,
	}

	// Act
	output, err := uc.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}
