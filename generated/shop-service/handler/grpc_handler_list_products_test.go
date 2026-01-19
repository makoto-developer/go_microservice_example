package handler_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/handler"
	pb "github.com/makoto-developer/go_microservice_example/generated/shop-service/proto/shop_service/v1"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/usecase"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockListProductsUsecase is a mock implementation of ListProductsUsecase
type MockListProductsUsecase struct {
	mock.Mock
}

func (m *MockListProductsUsecase) Execute(ctx context.Context, input usecase.ListProductsInput) (*usecase.ListProductsOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.ListProductsOutput), args.Error(1)
}

func TestShopServiceHandler_ListProducts_Success(t *testing.T) {
	// Arrange
	mockUC := new(MockListProductsUsecase)

	// 他のusecaseは実際には使われないのでnilでOK
	h := handler.NewShopServiceHandler(
		nil, // registerShopUsecase
		nil, // updateShopInfoUsecase
		nil, // toggleShopPublishUsecase
		nil, // registerProductUsecase
		nil, // updateProductUsecase
		nil, // deleteProductUsecase
		nil, // toggleProductPublishUsecase
		nil, // uploadProductImageUsecase
		nil, // getProductUsecase
		mockUC, // listProductsUsecase
		nil, // manageProductVariationUsecase
		nil, // listOrdersUsecase
		nil, // getOrderDetailUsecase
		nil, // updateOrderStatusUsecase
		nil, // getSalesReportUsecase
		nil, // exportSalesDataUsecase
	)

	ctx := context.Background()
	shopID := "11111111-1111-1111-1111-111111111111"
	req := &pb.ListProductsRequest{
		ShopId:        shopID,
		Category:      "オーディオ",
		PublishedOnly: true,
		Limit:         10,
		Offset:        0,
	}

	price, _ := decimal.NewFromString("29800")
	weight, _ := decimal.NewFromString("0.055")
	size := "5.4 x 4.6 x 2.1 cm"
	janCode := "4901234567890"

	expectedOutput := &usecase.ListProductsOutput{
		Products: []*domain.Product{
			{
				Id:            uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
				ShopId:        uuid.MustParse(shopID),
				Name:          "ワイヤレスイヤホン Pro",
				Description:   "高音質ワイヤレスイヤホン",
				Price:         price,
				Category:      "オーディオ",
				StockQuantity: 50,
				Weight:        &weight,
				Size:          &size,
				JanCode:       &janCode,
				Published:     true,
				Deleted:       false,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		},
	}

	mockUC.On("Execute", ctx, mock.MatchedBy(func(input usecase.ListProductsInput) bool {
		return input.ShopID == uuid.MustParse(shopID) &&
			input.Category == "オーディオ" &&
			input.PublishedOnly == true &&
			input.Limit == 10 &&
			input.Offset == 0
	})).Return(expectedOutput, nil)

	// Act
	resp, err := h.ListProducts(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(resp.Products))

	product := resp.Products[0]
	assert.Equal(t, "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", product.Id)
	assert.Equal(t, shopID, product.ShopId)
	assert.Equal(t, "ワイヤレスイヤホン Pro", product.Name)
	assert.Equal(t, "高音質ワイヤレスイヤホン", product.Description)
	assert.Equal(t, "29800", product.Price)
	assert.Equal(t, "オーディオ", product.Category)
	assert.Equal(t, int32(50), product.StockQuantity)
	assert.Equal(t, "0.055", product.Weight)
	assert.Equal(t, "5.4 x 4.6 x 2.1 cm", product.Size)
	assert.Equal(t, "4901234567890", product.JanCode)
	assert.True(t, product.Published)

	mockUC.AssertExpectations(t)
}

func TestShopServiceHandler_ListProducts_EmptyShopID(t *testing.T) {
	// Arrange
	mockUC := new(MockListProductsUsecase)
	h := handler.NewShopServiceHandler(
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		mockUC,
		nil, nil, nil, nil, nil, nil,
	)

	ctx := context.Background()
	req := &pb.ListProductsRequest{
		ShopId:        "",  // 空のshop_id (全ショップ対象)
		Category:      "",
		PublishedOnly: false,
		Limit:         100,
		Offset:        0,
	}

	price1, _ := decimal.NewFromString("29800")
	price2, _ := decimal.NewFromString("45000")

	expectedOutput := &usecase.ListProductsOutput{
		Products: []*domain.Product{
			{
				Id:            uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
				ShopId:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				Name:          "商品1",
				Price:         price1,
				Category:      "カテゴリ1",
				StockQuantity: 50,
				Published:     true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			{
				Id:            uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
				ShopId:        uuid.MustParse("22222222-2222-2222-2222-222222222222"),
				Name:          "商品2",
				Price:         price2,
				Category:      "カテゴリ2",
				StockQuantity: 30,
				Published:     true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		},
	}

	mockUC.On("Execute", ctx, mock.MatchedBy(func(input usecase.ListProductsInput) bool {
		return input.ShopID == uuid.Nil  // 空の場合はuuid.Nilになる
	})).Return(expectedOutput, nil)

	// Act
	resp, err := h.ListProducts(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Products))

	mockUC.AssertExpectations(t)
}

func TestShopServiceHandler_ListProducts_InvalidShopID(t *testing.T) {
	// Arrange
	mockUC := new(MockListProductsUsecase)
	h := handler.NewShopServiceHandler(
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		mockUC,
		nil, nil, nil, nil, nil, nil,
	)

	ctx := context.Background()
	req := &pb.ListProductsRequest{
		ShopId:        "invalid-uuid",  // 不正なUUID
		Category:      "",
		PublishedOnly: false,
		Limit:         10,
		Offset:        0,
	}

	// Act
	resp, err := h.ListProducts(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)

	// usecaseは呼ばれないはず
	mockUC.AssertNotCalled(t, "Execute")
}

func TestShopServiceHandler_ListProducts_UsecaseError(t *testing.T) {
	// Arrange
	mockUC := new(MockListProductsUsecase)
	h := handler.NewShopServiceHandler(
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		mockUC,
		nil, nil, nil, nil, nil, nil,
	)

	ctx := context.Background()
	shopID := "11111111-1111-1111-1111-111111111111"
	req := &pb.ListProductsRequest{
		ShopId:        shopID,
		Category:      "",
		PublishedOnly: false,
		Limit:         10,
		Offset:        0,
	}

	mockUC.On("Execute", ctx, mock.Anything).Return(nil, assert.AnError)

	// Act
	resp, err := h.ListProducts(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)

	mockUC.AssertExpectations(t)
}

func TestShopServiceHandler_ListProducts_EmptyResult(t *testing.T) {
	// Arrange
	mockUC := new(MockListProductsUsecase)
	h := handler.NewShopServiceHandler(
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		mockUC,
		nil, nil, nil, nil, nil, nil,
	)

	ctx := context.Background()
	req := &pb.ListProductsRequest{
		ShopId:        "11111111-1111-1111-1111-111111111111",
		Category:      "存在しないカテゴリ",
		PublishedOnly: true,
		Limit:         10,
		Offset:        0,
	}

	expectedOutput := &usecase.ListProductsOutput{
		Products: []*domain.Product{},
	}

	mockUC.On("Execute", ctx, mock.Anything).Return(expectedOutput, nil)

	// Act
	resp, err := h.ListProducts(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 0, len(resp.Products))

	mockUC.AssertExpectations(t)
}

func TestShopServiceHandler_ListProducts_WithOptionalFields(t *testing.T) {
	// Arrange
	mockUC := new(MockListProductsUsecase)
	h := handler.NewShopServiceHandler(
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		mockUC,
		nil, nil, nil, nil, nil, nil,
	)

	ctx := context.Background()
	req := &pb.ListProductsRequest{
		ShopId:        "11111111-1111-1111-1111-111111111111",
		Category:      "",
		PublishedOnly: false,
		Limit:         10,
		Offset:        0,
	}

	price, _ := decimal.NewFromString("29800")

	expectedOutput := &usecase.ListProductsOutput{
		Products: []*domain.Product{
			{
				Id:            uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
				ShopId:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				Name:          "シンプル商品",
				Description:   "必須フィールドのみ",
				Price:         price,
				Category:      "カテゴリ",
				StockQuantity: 10,
				Weight:        nil,  // オプションフィールドはnil
				Size:          nil,
				JanCode:       nil,
				Published:     true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		},
	}

	mockUC.On("Execute", ctx, mock.Anything).Return(expectedOutput, nil)

	// Act
	resp, err := h.ListProducts(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(resp.Products))

	product := resp.Products[0]
	assert.Equal(t, "", product.Weight)
	assert.Equal(t, "", product.Size)
	assert.Equal(t, "", product.JanCode)

	mockUC.AssertExpectations(t)
}

func TestShopServiceHandler_ListProducts_LargeDataset(t *testing.T) {
	// Arrange
	mockUC := new(MockListProductsUsecase)
	h := handler.NewShopServiceHandler(
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		mockUC,
		nil, nil, nil, nil, nil, nil,
	)

	ctx := context.Background()
	req := &pb.ListProductsRequest{
		ShopId:        "11111111-1111-1111-1111-111111111111",
		Category:      "",
		PublishedOnly: false,
		Limit:         100,
		Offset:        0,
	}

	// 100個の商品を作成
	products := make([]*domain.Product, 100)
	price, _ := decimal.NewFromString("1000")
	for i := 0; i < 100; i++ {
		products[i] = &domain.Product{
			Id:            uuid.New(),
			ShopId:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Name:          "商品" + string(rune(i)),
			Description:   "説明",
			Price:         price,
			Category:      "カテゴリ",
			StockQuantity: 10,
			Published:     true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
	}

	expectedOutput := &usecase.ListProductsOutput{
		Products: products,
	}

	mockUC.On("Execute", ctx, mock.Anything).Return(expectedOutput, nil)

	// Act
	resp, err := h.ListProducts(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 100, len(resp.Products))

	mockUC.AssertExpectations(t)
}
