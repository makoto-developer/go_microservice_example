package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	pb "github.com/makoto-developer/go_microservice_example/proto/shop_service/v1"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ShopServiceHandler struct {
	pb.UnimplementedShopServiceServer
	registerShopUsecase         usecase.RegisterShopUsecase
	updateShopUsecase           usecase.UpdateShopUsecase
	toggleShopPublishUsecase    usecase.ToggleShopPublishUsecase
	getShopUsecase              usecase.GetShopUsecase
	registerProductUsecase      usecase.RegisterProductUsecase
	updateProductUsecase        usecase.UpdateProductUsecase
	deleteProductUsecase        usecase.DeleteProductUsecase
	toggleProductPublishUsecase usecase.ToggleProductPublishUsecase
	getProductUsecase           usecase.GetProductUsecase
	listProductsUsecase         usecase.ListProductsUsecase
	uploadProductImageUsecase   usecase.UploadProductImageUsecase
	manageVariationUsecase      usecase.ManageVariationUsecase
	listOrdersUsecase           usecase.ListOrdersUsecase
	updateOrderStatusUsecase    usecase.UpdateOrderStatusUsecase
	getSalesReportUsecase       usecase.GetSalesReportUsecase
}

func NewShopServiceHandler(
	registerShopUsecase usecase.RegisterShopUsecase,
	updateShopUsecase usecase.UpdateShopUsecase,
	toggleShopPublishUsecase usecase.ToggleShopPublishUsecase,
	getShopUsecase usecase.GetShopUsecase,
	registerProductUsecase usecase.RegisterProductUsecase,
	updateProductUsecase usecase.UpdateProductUsecase,
	deleteProductUsecase usecase.DeleteProductUsecase,
	toggleProductPublishUsecase usecase.ToggleProductPublishUsecase,
	getProductUsecase usecase.GetProductUsecase,
	listProductsUsecase usecase.ListProductsUsecase,
	uploadProductImageUsecase usecase.UploadProductImageUsecase,
	manageVariationUsecase usecase.ManageVariationUsecase,
	listOrdersUsecase usecase.ListOrdersUsecase,
	updateOrderStatusUsecase usecase.UpdateOrderStatusUsecase,
	getSalesReportUsecase usecase.GetSalesReportUsecase,
) *ShopServiceHandler {
	return &ShopServiceHandler{
		registerShopUsecase:         registerShopUsecase,
		updateShopUsecase:           updateShopUsecase,
		toggleShopPublishUsecase:    toggleShopPublishUsecase,
		getShopUsecase:              getShopUsecase,
		registerProductUsecase:      registerProductUsecase,
		updateProductUsecase:        updateProductUsecase,
		deleteProductUsecase:        deleteProductUsecase,
		toggleProductPublishUsecase: toggleProductPublishUsecase,
		getProductUsecase:           getProductUsecase,
		listProductsUsecase:         listProductsUsecase,
		uploadProductImageUsecase:   uploadProductImageUsecase,
		manageVariationUsecase:      manageVariationUsecase,
		listOrdersUsecase:           listOrdersUsecase,
		updateOrderStatusUsecase:    updateOrderStatusUsecase,
		getSalesReportUsecase:       getSalesReportUsecase,
	}
}

func (h *ShopServiceHandler) RegisterShop(ctx context.Context, req *pb.RegisterShopRequest) (*pb.RegisterShopResponse, error) {
	ownerID, err := uuid.Parse(req.OwnerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid owner ID")
	}

	var logoURL *string
	if req.LogoUrl != "" {
		logoURL = &req.LogoUrl
	}

	input := usecase.RegisterShopInput{
		OwnerID:       ownerID,
		Name:          req.Name,
		Description:   req.Description,
		LogoURL:       logoURL,
		OwnerName:     req.OwnerName,
		PhoneNumber:   req.PhoneNumber,
		BusinessHours: req.BusinessHours,
		ReturnPolicy:  req.ReturnPolicy,
		Categories:    req.Categories,
	}

	output, err := h.registerShopUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.RegisterShopResponse{
		ShopId:  output.ShopID.String(),
		Status:  domainStatusToProto(output.Status),
		Message: output.Message,
	}, nil
}

func (h *ShopServiceHandler) UpdateShop(ctx context.Context, req *pb.UpdateShopRequest) (*pb.UpdateShopResponse, error) {
	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid shop ID")
	}

	var logoURL *string
	if req.LogoUrl != "" {
		logoURL = &req.LogoUrl
	}

	input := usecase.UpdateShopInput{
		ShopID:        shopID,
		Name:          req.Name,
		Description:   req.Description,
		LogoURL:       logoURL,
		BusinessHours: req.BusinessHours,
		ReturnPolicy:  req.ReturnPolicy,
	}

	output, err := h.updateShopUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.UpdateShopResponse{
		ShopId:             output.ShopID.String(),
		RequiresReapproval: output.RequiresReapproval,
	}, nil
}

func (h *ShopServiceHandler) ToggleShopPublish(ctx context.Context, req *pb.ToggleShopPublishRequest) (*pb.ToggleShopPublishResponse, error) {
	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid shop ID")
	}

	input := usecase.ToggleShopPublishInput{
		ShopID:    shopID,
		Published: req.Published,
	}

	output, err := h.toggleShopPublishUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.ToggleShopPublishResponse{
		ShopId:    output.ShopID.String(),
		Published: output.Published,
	}, nil
}

func (h *ShopServiceHandler) GetShop(ctx context.Context, req *pb.GetShopRequest) (*pb.GetShopResponse, error) {
	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid shop ID")
	}

	input := usecase.GetShopInput{ShopID: shopID}
	output, err := h.getShopUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	var categories []string
	for _, cat := range output.Categories {
		categories = append(categories, cat.CategoryName)
	}

	shop := &pb.Shop{
		Id:            output.Shop.ID.String(),
		OwnerId:       output.Shop.OwnerID.String(),
		Name:          output.Shop.Name,
		Description:   output.Shop.Description,
		OwnerName:     output.Shop.OwnerName,
		PhoneNumber:   output.Shop.PhoneNumber,
		BusinessHours: output.Shop.BusinessHours,
		ReturnPolicy:  output.Shop.ReturnPolicy,
		Status:        domainStatusToProto(output.Shop.Status),
		Published:     output.Shop.Published,
		CreatedAt:     timestamppb.New(output.Shop.CreatedAt),
		UpdatedAt:     timestamppb.New(output.Shop.UpdatedAt),
	}
	if output.Shop.LogoURL != nil {
		shop.LogoUrl = *output.Shop.LogoURL
	}

	return &pb.GetShopResponse{Shop: shop}, nil
}

func (h *ShopServiceHandler) RegisterProduct(ctx context.Context, req *pb.RegisterProductRequest) (*pb.RegisterProductResponse, error) {
	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid shop ID")
	}

	price, err := parsePrice(req.Price)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid price")
	}

	input := usecase.RegisterProductInput{
		ShopID:        shopID,
		Name:          req.Name,
		Description:   req.Description,
		Price:         price,
		Category:      req.Category,
		StockQuantity: int(req.StockQuantity),
		Tags:          req.Tags,
	}

	if req.Weight != "" {
		weight, err := parseWeight(req.Weight)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid weight")
		}
		input.Weight = weight
	}
	if req.Size != "" {
		input.Size = &req.Size
	}
	if req.JanCode != "" {
		input.JANCode = &req.JanCode
	}

	output, err := h.registerProductUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.RegisterProductResponse{ProductId: output.ProductID.String()}, nil
}

func (h *ShopServiceHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	price, err := parsePrice(req.Price)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid price")
	}

	input := usecase.UpdateProductInput{
		ProductID:     productID,
		Name:          req.Name,
		Description:   req.Description,
		Price:         price,
		Category:      req.Category,
		StockQuantity: int(req.StockQuantity),
	}

	if req.Weight != "" {
		weight, err := parseWeight(req.Weight)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid weight")
		}
		input.Weight = weight
	}
	if req.Size != "" {
		input.Size = &req.Size
	}
	if req.JanCode != "" {
		input.JANCode = &req.JanCode
	}

	output, err := h.updateProductUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.UpdateProductResponse{ProductId: output.ProductID.String()}, nil
}

func (h *ShopServiceHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	input := usecase.DeleteProductInput{ProductID: productID}
	output, err := h.deleteProductUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.DeleteProductResponse{
		ProductId: output.ProductID.String(),
		Deleted:   output.Deleted,
	}, nil
}

func (h *ShopServiceHandler) ToggleProductPublish(ctx context.Context, req *pb.ToggleProductPublishRequest) (*pb.ToggleProductPublishResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	input := usecase.ToggleProductPublishInput{
		ProductID: productID,
		Published: req.Published,
	}

	output, err := h.toggleProductPublishUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.ToggleProductPublishResponse{
		ProductId: output.ProductID.String(),
		Published: output.Published,
	}, nil
}

func (h *ShopServiceHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	input := usecase.GetProductInput{ProductID: productID}
	output, err := h.getProductUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	product := &pb.Product{
		Id:            output.Product.ID.String(),
		ShopId:        output.Product.ShopID.String(),
		Name:          output.Product.Name,
		Description:   output.Product.Description,
		Price:         formatPrice(output.Product.Price),
		Category:      output.Product.Category,
		StockQuantity: int32(output.Product.StockQuantity),
		Weight:        formatWeight(output.Product.Weight),
		Published:     output.Product.Published,
		Deleted:       output.Product.Deleted,
		CreatedAt:     timestamppb.New(output.Product.CreatedAt),
		UpdatedAt:     timestamppb.New(output.Product.UpdatedAt),
	}
	if output.Product.Size != nil {
		product.Size = *output.Product.Size
	}
	if output.Product.JANCode != nil {
		product.JanCode = *output.Product.JANCode
	}

	return &pb.GetProductResponse{Product: product}, nil
}

func (h *ShopServiceHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid shop ID")
	}

	input := usecase.ListProductsInput{
		ShopID:         shopID,
		IncludeDeleted: false,
	}

	output, err := h.listProductsUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	var products []*pb.Product
	for _, p := range output.Products {
		product := &pb.Product{
			Id:            p.ID.String(),
			ShopId:        p.ShopID.String(),
			Name:          p.Name,
			Description:   p.Description,
			Price:         formatPrice(p.Price),
			Category:      p.Category,
			StockQuantity: int32(p.StockQuantity),
			Weight:        formatWeight(p.Weight),
			Published:     p.Published,
			Deleted:       p.Deleted,
			CreatedAt:     timestamppb.New(p.CreatedAt),
			UpdatedAt:     timestamppb.New(p.UpdatedAt),
		}
		if p.Size != nil {
			product.Size = *p.Size
		}
		if p.JANCode != nil {
			product.JanCode = *p.JANCode
		}
		products = append(products, product)
	}

	return &pb.ListProductsResponse{Products: products}, nil
}

func (h *ShopServiceHandler) UploadProductImage(ctx context.Context, req *pb.UploadProductImageRequest) (*pb.UploadProductImageResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	input := usecase.UploadProductImageInput{
		ProductID:    productID,
		ImageData:    req.ImageData,
		DisplayOrder: int(req.DisplayOrder),
	}

	output, err := h.uploadProductImageUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.UploadProductImageResponse{
		ImageId:    output.ImageID.String(),
		Url:        output.URL,
		Thumbnails: output.Thumbnails,
	}, nil
}

func (h *ShopServiceHandler) ManageVariation(ctx context.Context, req *pb.ManageVariationRequest) (*pb.ManageVariationResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	var variations []usecase.ProductVariationInput
	for _, v := range req.Variations {
		price, err := parsePrice(v.Price)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid variation price")
		}
		variations = append(variations, usecase.ProductVariationInput{
			SKU:            v.Sku,
			AttributeName:  v.AttributeName,
			AttributeValue: v.AttributeValue,
			Price:          price,
			StockQuantity:  int(v.StockQuantity),
		})
	}

	input := usecase.ManageVariationInput{
		ProductID:  productID,
		Variations: variations,
	}

	output, err := h.manageVariationUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	var variationIDs []string
	for _, id := range output.VariationIDs {
		variationIDs = append(variationIDs, id.String())
	}

	return &pb.ManageVariationResponse{VariationIds: variationIDs}, nil
}

func (h *ShopServiceHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid shop ID")
	}

	filter := repository.OrderFilter{
		ShopID:    shopID,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	if req.Status != pb.OrderStatus_ORDER_STATUS_UNSPECIFIED {
		orderStatus := domain.OrderStatus(req.Status.String())
		filter.Status = &orderStatus
	}

	if req.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", req.DateFrom)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid date_from format")
		}
		filter.DateFrom = &dateFrom
	}

	if req.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", req.DateTo)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid date_to format")
		}
		filter.DateTo = &dateTo
	}

	if req.CustomerName != "" {
		filter.CustomerName = &req.CustomerName
	}

	if req.ProductName != "" {
		filter.ProductName = &req.ProductName
	}

	input := usecase.ListOrdersInput{Filter: filter}
	output, err := h.listOrdersUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	var orders []*pb.OrderSummary
	for _, o := range output.Orders {
		order := &pb.OrderSummary{
			Id:          o.ID.String(),
			OrderNumber: o.OrderNumber,
			Status:      domainOrderStatusToProto(o.Status),
			TotalAmount: formatPrice(o.TotalAmount),
			CreatedAt:   timestamppb.New(o.CreatedAt),
		}
		orders = append(orders, order)
	}

	return &pb.ListOrdersResponse{
		Orders:     orders,
		TotalCount: int32(output.TotalCount),
	}, nil
}

func (h *ShopServiceHandler) GetOrderDetail(ctx context.Context, req *pb.GetOrderDetailRequest) (*pb.GetOrderDetailResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (h *ShopServiceHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid shop ID")
	}

	input := usecase.UpdateOrderStatusInput{
		OrderID:   orderID,
		ShopID:    shopID,
		NewStatus: domain.OrderStatus(req.NewStatus.String()),
	}

	if req.TrackingNumber != "" {
		input.TrackingNumber = &req.TrackingNumber
	}

	if req.Carrier != pb.Carrier_CARRIER_UNSPECIFIED {
		input.Carrier = protoCarrierToDomain(req.Carrier)
	}

	output, err := h.updateOrderStatusUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.UpdateOrderStatusResponse{
		OrderId: output.OrderID.String(),
		Status:  domainOrderStatusToProto(output.Status),
	}, nil
}

func (h *ShopServiceHandler) GetSalesReport(ctx context.Context, req *pb.GetSalesReportRequest) (*pb.GetSalesReportResponse, error) {
	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid shop ID")
	}

	dateFrom, err := time.Parse("2006-01-02", req.DateFrom)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date_from format")
	}
	dateTo, err := time.Parse("2006-01-02", req.DateTo)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date_to format")
	}

	input := usecase.GetSalesReportInput{
		ShopID:     shopID,
		ReportType: req.ReportType,
		DateFrom:   dateFrom,
		DateTo:     dateTo,
	}

	output, err := h.getSalesReportUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	var reportData []*pb.SalesData
	for _, data := range output.ReportData {
		reportData = append(reportData, &pb.SalesData{
			Date:              data.Date.Format("2006-01-02"),
			TotalSales:        formatPrice(data.TotalSales),
			OrderCount:        int32(data.OrderCount),
			AverageOrderValue: formatPrice(data.AverageOrderValue),
		})
	}

	return &pb.GetSalesReportResponse{
		ReportData: reportData,
		Summary: &pb.SalesSummary{
			TotalSales:        formatPrice(output.Summary.TotalSales),
			TotalOrders:       int32(output.Summary.TotalOrders),
			AverageOrderValue: formatPrice(output.Summary.AverageOrderValue),
		},
	}, nil
}

func (h *ShopServiceHandler) ExportSalesData(ctx context.Context, req *pb.ExportSalesDataRequest) (*pb.ExportSalesDataResponse, error) {
	return &pb.ExportSalesDataResponse{
		CsvUrl:    "https://example.com/exports/sales-" + req.ShopId + ".csv",
		ExpiresAt: timestamppb.New(time.Now().Add(24 * time.Hour)),
	}, nil
}

func mapDomainError(err error) error {
	switch err {
	case domain.ErrShopNotFound, domain.ErrProductNotFound, domain.ErrOrderNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrShopAlreadyExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case domain.ErrShopNotApproved, domain.ErrUnauthorizedAccess:
		return status.Error(codes.PermissionDenied, err.Error())
	case domain.ErrInvalidShopData, domain.ErrInvalidProductData, domain.ErrInvalidStatusTransition, domain.ErrInvalidDateRange:
		return status.Error(codes.InvalidArgument, err.Error())
	case domain.ErrInsufficientStock:
		return status.Error(codes.FailedPrecondition, err.Error())
	case domain.ErrMaxImagesExceeded, domain.ErrImageTooLarge, domain.ErrInvalidImageFormat:
		return status.Error(codes.InvalidArgument, err.Error())
	case domain.ErrDuplicateSKU:
		return status.Error(codes.AlreadyExists, err.Error())
	case domain.ErrNoDataFound:
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
