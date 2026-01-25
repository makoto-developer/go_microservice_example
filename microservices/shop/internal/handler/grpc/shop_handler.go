package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	pb "github.com/makoto-developer/go_microservice_example/generated/shop/proto/shop_service/v1"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShopServiceHandler struct {
	pb.UnimplementedShopServiceServer
	shopRegistrationUsecase  usecase.ShopRegistrationUsecase
	productManagementUsecase usecase.ProductManagementUsecase
	shopRepo                 repository.ShopRepository
}

func NewShopServiceHandler(
	shopRegistrationUsecase usecase.ShopRegistrationUsecase,
	productManagementUsecase usecase.ProductManagementUsecase,
	shopRepo repository.ShopRepository,
) *ShopServiceHandler {
	return &ShopServiceHandler{
		shopRegistrationUsecase:  shopRegistrationUsecase,
		productManagementUsecase: productManagementUsecase,
		shopRepo:                 shopRepo,
	}
}

// RegisterShop registers a new shop
func (h *ShopServiceHandler) RegisterShop(ctx context.Context, req *pb.RegisterShopRequest) (*pb.RegisterShopResponse, error) {
	ownerID, err := uuid.Parse(req.OwnerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid owner_id: %v", err)
	}

	categoryIDs := make([]uuid.UUID, 0, len(req.Categories))
	for _, catStr := range req.Categories {
		catID, err := uuid.Parse(catStr)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid category ID: %v", err)
		}
		categoryIDs = append(categoryIDs, catID)
	}

	input := usecase.ShopRegistrationInput{
		OwnerID:       ownerID,
		Name:          req.Name,
		Description:   req.Description,
		LogoURL:       req.LogoUrl,
		OwnerName:     req.OwnerName,
		PhoneNumber:   req.PhoneNumber,
		BusinessHours: req.BusinessHours,
		ReturnPolicy:  req.ReturnPolicy,
		CategoryIDs:   categoryIDs,
	}

	output, err := h.shopRegistrationUsecase.Execute(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register shop: %v", err)
	}

	return &pb.RegisterShopResponse{
		ShopId: output.ShopID.String(),
		Status: convertToProtoShopStatus(output.Status),
		Message: "Shop registered successfully and pending approval",
	}, nil
}

// UpdateShop updates shop information
func (h *ShopServiceHandler) UpdateShop(ctx context.Context, req *pb.UpdateShopRequest) (*pb.UpdateShopResponse, error) {
	// TODO: Implement shop update logic
	return &pb.UpdateShopResponse{
		ShopId: req.ShopId,
		RequiresReapproval: false,
	}, nil
}

// ToggleShopPublish toggles shop publish status
func (h *ShopServiceHandler) ToggleShopPublish(ctx context.Context, req *pb.ToggleShopPublishRequest) (*pb.ToggleShopPublishResponse, error) {
	// TODO: Implement toggle publish logic
	return &pb.ToggleShopPublishResponse{
		ShopId: req.ShopId,
		Published: req.Published,
	}, nil
}

// GetShop retrieves shop information
func (h *ShopServiceHandler) GetShop(ctx context.Context, req *pb.GetShopRequest) (*pb.GetShopResponse, error) {
	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid shop_id: %v", err)
	}

	shop, err := h.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shop not found: %v", err)
	}

	return &pb.GetShopResponse{
		Shop: convertToProtoShop(shop),
	}, nil
}

// ListShops lists all public shops
func (h *ShopServiceHandler) ListShops(ctx context.Context, req *pb.ListShopsRequest) (*pb.ListShopsResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 100
	}
	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	shops, err := h.shopRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list shops: %v", err)
	}

	// 公開ショップのみフィルタリング
	protoShops := make([]*pb.Shop, 0, len(shops))
	for _, shop := range shops {
		if req.PublishedOnly && !shop.Published {
			continue
		}
		protoShops = append(protoShops, convertToProtoShop(shop))
	}

	return &pb.ListShopsResponse{
		Shops:      protoShops,
		TotalCount: int32(len(protoShops)),
	}, nil
}

// GetShopsByOwner retrieves shops by owner ID
func (h *ShopServiceHandler) GetShopsByOwner(ctx context.Context, req *pb.GetShopsByOwnerRequest) (*pb.ListShopsResponse, error) {
	ownerID, err := uuid.Parse(req.OwnerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid owner_id: %v", err)
	}

	shops, err := h.shopRepo.GetByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get shops: %v", err)
	}

	protoShops := make([]*pb.Shop, 0, len(shops))
	for _, shop := range shops {
		protoShops = append(protoShops, convertToProtoShop(shop))
	}

	return &pb.ListShopsResponse{
		Shops:      protoShops,
		TotalCount: int32(len(protoShops)),
	}, nil
}

// RegisterProduct registers a new product
func (h *ShopServiceHandler) RegisterProduct(ctx context.Context, req *pb.RegisterProductRequest) (*pb.RegisterProductResponse, error) {
	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid shop_id: %v", err)
	}

	input := usecase.ProductCreateInput{
		ShopID:        shopID,
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		Category:      req.Category,
		StockQuantity: int(req.StockQuantity),
	}

	if req.Weight != "" {
		var weight float64
		_, err = fmt.Sscanf(req.Weight, "%f", &weight)
		if err == nil {
			input.Weight = &weight
		}
	}

	if req.Size != "" {
		input.Size = &req.Size
	}

	if req.JanCode != "" {
		input.JANCode = &req.JanCode
	}

	productID, err := h.productManagementUsecase.CreateProduct(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &pb.RegisterProductResponse{
		ProductId: productID.String(),
	}, nil
}

// UpdateProduct updates product information
func (h *ShopServiceHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product_id: %v", err)
	}

	input := usecase.ProductUpdateInput{
		ID:            productID,
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		Category:      req.Category,
		StockQuantity: int(req.StockQuantity),
	}

	if req.Weight != "" {
		var weight float64
		_, err = fmt.Sscanf(req.Weight, "%f", &weight)
		if err == nil {
			input.Weight = &weight
		}
	}

	if req.Size != "" {
		input.Size = &req.Size
	}

	if req.JanCode != "" {
		input.JANCode = &req.JanCode
	}

	err = h.productManagementUsecase.UpdateProduct(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	return &pb.UpdateProductResponse{
		ProductId: productID.String(),
	}, nil
}

// DeleteProduct deletes a product
func (h *ShopServiceHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product_id: %v", err)
	}

	err = h.productManagementUsecase.DeleteProduct(ctx, productID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	return &pb.DeleteProductResponse{
		ProductId: productID.String(),
		Deleted: true,
	}, nil
}

// ToggleProductPublish toggles product publish status
func (h *ShopServiceHandler) ToggleProductPublish(ctx context.Context, req *pb.ToggleProductPublishRequest) (*pb.ToggleProductPublishResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product_id: %v", err)
	}

	if req.Published {
		err = h.productManagementUsecase.PublishProduct(ctx, productID)
	} else {
		err = h.productManagementUsecase.UnpublishProduct(ctx, productID)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to toggle product publish: %v", err)
	}

	return &pb.ToggleProductPublishResponse{
		ProductId: productID.String(),
		Published: req.Published,
	}, nil
}

// GetProduct retrieves product information
func (h *ShopServiceHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product_id: %v", err)
	}

	product, err := h.productManagementUsecase.GetProduct(ctx, productID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
	}

	return &pb.GetProductResponse{
		Product: convertToProtoProduct(product),
	}, nil
}

// ListProducts lists products
func (h *ShopServiceHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	// shop_idが空の場合は空のリストを返す（将来的には全商品を返す実装に変更可能）
	if req.ShopId == "" {
		return &pb.ListProductsResponse{
			Products:   []*pb.Product{},
			TotalCount: 0,
		}, nil
	}

	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid shop_id: %v", err)
	}

	products, err := h.productManagementUsecase.ListProductsByShop(ctx, shopID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}

	protoProducts := make([]*pb.Product, 0, len(products))
	for _, p := range products {
		protoProducts = append(protoProducts, convertToProtoProduct(p))
	}

	return &pb.ListProductsResponse{
		Products:   protoProducts,
		TotalCount: int32(len(protoProducts)),
	}, nil
}

// UploadProductImage uploads a product image
func (h *ShopServiceHandler) UploadProductImage(ctx context.Context, req *pb.UploadProductImageRequest) (*pb.UploadProductImageResponse, error) {
	// TODO: Implement image upload logic
	return &pb.UploadProductImageResponse{
		ImageId: uuid.New().String(),
		Url: "https://example.com/image.jpg",
		Thumbnails: map[string]string{
			"200": "https://example.com/thumb_200.jpg",
			"400": "https://example.com/thumb_400.jpg",
			"800": "https://example.com/thumb_800.jpg",
		},
	}, nil
}

// ManageVariation manages product variations
func (h *ShopServiceHandler) ManageVariation(ctx context.Context, req *pb.ManageVariationRequest) (*pb.ManageVariationResponse, error) {
	// TODO: Implement variation management logic
	return &pb.ManageVariationResponse{
		VariationIds: []string{uuid.New().String()},
	}, nil
}

// ListOrders lists orders for a shop
func (h *ShopServiceHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	// TODO: Implement order listing logic
	return &pb.ListOrdersResponse{
		Orders: []*pb.OrderSummary{},
		TotalCount: 0,
	}, nil
}

// GetOrderDetail retrieves order detail
func (h *ShopServiceHandler) GetOrderDetail(ctx context.Context, req *pb.GetOrderDetailRequest) (*pb.GetOrderDetailResponse, error) {
	// TODO: Implement order detail retrieval logic
	return &pb.GetOrderDetailResponse{}, nil
}

// UpdateOrderStatus updates order status
func (h *ShopServiceHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	// TODO: Implement order status update logic
	return &pb.UpdateOrderStatusResponse{
		OrderId: req.OrderId,
		Status: req.NewStatus,
	}, nil
}

// GetSalesReport retrieves sales report
func (h *ShopServiceHandler) GetSalesReport(ctx context.Context, req *pb.GetSalesReportRequest) (*pb.GetSalesReportResponse, error) {
	// TODO: Implement sales report logic
	return &pb.GetSalesReportResponse{
		ReportData: []*pb.SalesData{},
		Summary: &pb.SalesSummary{
			TotalSales: "0",
			TotalOrders: 0,
			AverageOrderValue: "0",
		},
	}, nil
}

// ExportSalesData exports sales data
func (h *ShopServiceHandler) ExportSalesData(ctx context.Context, req *pb.ExportSalesDataRequest) (*pb.ExportSalesDataResponse, error) {
	// TODO: Implement sales data export logic
	return &pb.ExportSalesDataResponse{
		CsvUrl: "https://example.com/export.csv",
	}, nil
}
