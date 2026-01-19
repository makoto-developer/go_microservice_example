package handler_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/handler"
	pb "github.com/makoto-developer/go_microservice_example/generated/shop-service/proto/shop_service/v1"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/usecase"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGetProductUsecase for testing
type MockGetProductUsecase struct {
	mock.Mock
}

func (m *MockGetProductUsecase) Execute(ctx context.Context, input usecase.GetProductInput) (*usecase.GetProductOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.GetProductOutput), args.Error(1)
}

func TestGetProductHandler_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUsecase := new(MockGetProductUsecase)

	productID := uuid.New()
	shopID := uuid.New()
	price := decimal.NewFromInt(29800)
	weight := decimal.NewFromFloat(0.055)
	size := "5.4 x 4.6 x 2.1 cm"
	janCode := "4901234567890"

	product := &domain.Product{
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

	mockUsecase.On("Execute", ctx, mock.MatchedBy(func(input usecase.GetProductInput) bool {
		return input.ProductID == productID
	})).Return(&usecase.GetProductOutput{
		Product: product,
	}, nil)

	h := handler.NewShopServiceHandler(
		nil, nil, nil, nil, nil, nil, nil, nil,
		mockUsecase, // get_productUsecase
		nil, nil, nil, nil, nil, nil, nil,
	)

	req := &pb.GetProductRequest{
		ProductId: productID.String(),
	}

	// Act
	resp, err := h.GetProduct(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Product)
	assert.Equal(t, productID.String(), resp.Product.Id)
	assert.Equal(t, shopID.String(), resp.Product.ShopId)
	assert.Equal(t, "ワイヤレスイヤホン Pro", resp.Product.Name)
	assert.Equal(t, "ノイズキャンセリング機能付きの高音質ワイヤレスイヤホン", resp.Product.Description)
	assert.Equal(t, "29800", resp.Product.Price)
	assert.Equal(t, "オーディオ", resp.Product.Category)
	assert.Equal(t, int32(50), resp.Product.StockQuantity)
	assert.Equal(t, "0.055", resp.Product.Weight)
	assert.Equal(t, "5.4 x 4.6 x 2.1 cm", resp.Product.Size)
	assert.Equal(t, "4901234567890", resp.Product.JanCode)
	assert.True(t, resp.Product.Published)
	mockUsecase.AssertExpectations(t)
}

func TestGetProductHandler_InvalidProductID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUsecase := new(MockGetProductUsecase)

	h := handler.NewShopServiceHandler(
		nil, nil, nil, nil, nil, nil, nil, nil,
		mockUsecase, // get_productUsecase
		nil, nil, nil, nil, nil, nil, nil,
	)

	req := &pb.GetProductRequest{
		ProductId: "invalid-uuid",
	}

	// Act
	resp, err := h.GetProduct(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid product_id")
	mockUsecase.AssertExpectations(t)
}

func TestGetProductHandler_ProductNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUsecase := new(MockGetProductUsecase)

	productID := uuid.New()

	mockUsecase.On("Execute", ctx, mock.MatchedBy(func(input usecase.GetProductInput) bool {
		return input.ProductID == productID
	})).Return(nil, fmt.Errorf("failed to find product"))

	h := handler.NewShopServiceHandler(
		nil, nil, nil, nil, nil, nil, nil, nil,
		mockUsecase, // get_productUsecase
		nil, nil, nil, nil, nil, nil, nil,
	)

	req := &pb.GetProductRequest{
		ProductId: productID.String(),
	}

	// Act
	resp, err := h.GetProduct(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to get product")
	mockUsecase.AssertExpectations(t)
}

func TestGetProductHandler_WithoutOptionalFields(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUsecase := new(MockGetProductUsecase)

	productID := uuid.New()
	shopID := uuid.New()
	price := decimal.NewFromInt(10000)

	product := &domain.Product{
		Id:            productID,
		ShopId:        shopID,
		Name:          "シンプル商品",
		Description:   "オプションフィールドなし",
		Price:         price,
		Category:      "その他",
		StockQuantity: 10,
		Weight:        nil, // オプションフィールドなし
		Size:          nil,
		JanCode:       nil,
		Published:     false,
		Deleted:       false,
	}

	mockUsecase.On("Execute", ctx, mock.MatchedBy(func(input usecase.GetProductInput) bool {
		return input.ProductID == productID
	})).Return(&usecase.GetProductOutput{
		Product: product,
	}, nil)

	h := handler.NewShopServiceHandler(
		nil, nil, nil, nil, nil, nil, nil, nil,
		mockUsecase, // get_productUsecase
		nil, nil, nil, nil, nil, nil, nil,
	)

	req := &pb.GetProductRequest{
		ProductId: productID.String(),
	}

	// Act
	resp, err := h.GetProduct(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Product)
	assert.Equal(t, productID.String(), resp.Product.Id)
	assert.Equal(t, "シンプル商品", resp.Product.Name)
	assert.Equal(t, "", resp.Product.Weight)  // オプションフィールドは空文字列
	assert.Equal(t, "", resp.Product.Size)
	assert.Equal(t, "", resp.Product.JanCode)
	assert.False(t, resp.Product.Published)
	mockUsecase.AssertExpectations(t)
}
