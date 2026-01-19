package handler

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	pb "github.com/makoto-developer/go_microservice_example/generated/shop-service/proto/shop_service/v1"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/usecase"
)

// ShopServiceHandler implements gRPC handler
type ShopServiceHandler struct {
	pb.UnimplementedShopServiceServer
	register_shopUsecase usecase.RegisterShopUsecase
	update_shop_infoUsecase usecase.UpdateShopInfoUsecase
	toggle_shop_publishUsecase usecase.ToggleShopPublishUsecase
	register_productUsecase usecase.RegisterProductUsecase
	update_productUsecase usecase.UpdateProductUsecase
	delete_productUsecase usecase.DeleteProductUsecase
	toggle_product_publishUsecase usecase.ToggleProductPublishUsecase
	upload_product_imageUsecase usecase.UploadProductImageUsecase
	get_productUsecase usecase.GetProductUsecase
	list_productsUsecase usecase.ListProductsUsecase
	manage_product_variationUsecase usecase.ManageProductVariationUsecase
	list_ordersUsecase usecase.ListOrdersUsecase
	get_order_detailUsecase usecase.GetOrderDetailUsecase
	update_order_statusUsecase usecase.UpdateOrderStatusUsecase
	get_sales_reportUsecase usecase.GetSalesReportUsecase
	export_sales_dataUsecase usecase.ExportSalesDataUsecase
}

// NewShopServiceHandler creates a new handler instance
func NewShopServiceHandler(
	register_shopUsecase usecase.RegisterShopUsecase,
	update_shop_infoUsecase usecase.UpdateShopInfoUsecase,
	toggle_shop_publishUsecase usecase.ToggleShopPublishUsecase,
	register_productUsecase usecase.RegisterProductUsecase,
	update_productUsecase usecase.UpdateProductUsecase,
	delete_productUsecase usecase.DeleteProductUsecase,
	toggle_product_publishUsecase usecase.ToggleProductPublishUsecase,
	upload_product_imageUsecase usecase.UploadProductImageUsecase,
	get_productUsecase usecase.GetProductUsecase,
	list_productsUsecase usecase.ListProductsUsecase,
	manage_product_variationUsecase usecase.ManageProductVariationUsecase,
	list_ordersUsecase usecase.ListOrdersUsecase,
	get_order_detailUsecase usecase.GetOrderDetailUsecase,
	update_order_statusUsecase usecase.UpdateOrderStatusUsecase,
	get_sales_reportUsecase usecase.GetSalesReportUsecase,
	export_sales_dataUsecase usecase.ExportSalesDataUsecase,
) *ShopServiceHandler {
	return &ShopServiceHandler{
		register_shopUsecase: register_shopUsecase,
		update_shop_infoUsecase: update_shop_infoUsecase,
		toggle_shop_publishUsecase: toggle_shop_publishUsecase,
		register_productUsecase: register_productUsecase,
		update_productUsecase: update_productUsecase,
		delete_productUsecase: delete_productUsecase,
		toggle_product_publishUsecase: toggle_product_publishUsecase,
		upload_product_imageUsecase: upload_product_imageUsecase,
		get_productUsecase: get_productUsecase,
		list_productsUsecase: list_productsUsecase,
		manage_product_variationUsecase: manage_product_variationUsecase,
		list_ordersUsecase: list_ordersUsecase,
		get_order_detailUsecase: get_order_detailUsecase,
		update_order_statusUsecase: update_order_statusUsecase,
		get_sales_reportUsecase: get_sales_reportUsecase,
		export_sales_dataUsecase: export_sales_dataUsecase,
	}
}

// RegisterShop handles RegisterShop RPC
func (h *ShopServiceHandler) RegisterShop(
	ctx context.Context,
	req *pb.RegisterShopRequest,
) (*pb.RegisterShopResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.RegisterShopResponse{}, nil
}

// UpdateShop handles UpdateShop RPC
func (h *ShopServiceHandler) UpdateShop(
	ctx context.Context,
	req *pb.UpdateShopRequest,
) (*pb.UpdateShopResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.UpdateShopResponse{}, nil
}

// ToggleShopPublish handles ToggleShopPublish RPC
func (h *ShopServiceHandler) ToggleShopPublish(
	ctx context.Context,
	req *pb.ToggleShopPublishRequest,
) (*pb.ToggleShopPublishResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.ToggleShopPublishResponse{}, nil
}

// GetShop handles GetShop RPC
func (h *ShopServiceHandler) GetShop(
	ctx context.Context,
	req *pb.GetShopRequest,
) (*pb.GetShopResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.GetShopResponse{}, nil
}

// RegisterProduct handles RegisterProduct RPC
func (h *ShopServiceHandler) RegisterProduct(
	ctx context.Context,
	req *pb.RegisterProductRequest,
) (*pb.RegisterProductResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.RegisterProductResponse{}, nil
}

// UpdateProduct handles UpdateProduct RPC
func (h *ShopServiceHandler) UpdateProduct(
	ctx context.Context,
	req *pb.UpdateProductRequest,
) (*pb.UpdateProductResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.UpdateProductResponse{}, nil
}

// DeleteProduct handles DeleteProduct RPC
func (h *ShopServiceHandler) DeleteProduct(
	ctx context.Context,
	req *pb.DeleteProductRequest,
) (*pb.DeleteProductResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.DeleteProductResponse{}, nil
}

// ToggleProductPublish handles ToggleProductPublish RPC
func (h *ShopServiceHandler) ToggleProductPublish(
	ctx context.Context,
	req *pb.ToggleProductPublishRequest,
) (*pb.ToggleProductPublishResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.ToggleProductPublishResponse{}, nil
}

// GetProduct handles GetProduct RPC
func (h *ShopServiceHandler) GetProduct(
	ctx context.Context,
	req *pb.GetProductRequest,
) (*pb.GetProductResponse, error) {
	log.Printf("GetProduct called: product_id=%s", req.ProductId)

	// Parse product ID
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		log.Printf("Invalid product_id: %v", err)
		return nil, fmt.Errorf("invalid product_id: %w", err)
	}

	// Execute usecase
	output, err := h.get_productUsecase.Execute(ctx, usecase.GetProductInput{
		ProductID: productID,
	})
	if err != nil {
		log.Printf("GetProduct usecase error: %v", err)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Convert to proto
	product := output.Product
	pbProduct := &pb.Product{
		Id:            product.Id.String(),
		ShopId:        product.ShopId.String(),
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price.String(),
		Category:      product.Category,
		StockQuantity: int32(product.StockQuantity),
		Published:     product.Published,
	}

	if product.Weight != nil {
		pbProduct.Weight = product.Weight.String()
	}
	if product.Size != nil {
		pbProduct.Size = *product.Size
	}
	if product.JanCode != nil {
		pbProduct.JanCode = *product.JanCode
	}

	log.Printf("GetProduct returning product: %s", product.Name)
	return &pb.GetProductResponse{
		Product: pbProduct,
	}, nil
}

// ListProducts handles ListProducts RPC
func (h *ShopServiceHandler) ListProducts(
	ctx context.Context,
	req *pb.ListProductsRequest,
) (*pb.ListProductsResponse, error) {
	fmt.Fprintf(os.Stderr, "!!! ListProducts CALLED !!!\n")
	log.Printf("ListProducts called: shop_id=%s, category=%s, published_only=%v, limit=%d, offset=%d",
		req.ShopId, req.Category, req.PublishedOnly, req.Limit, req.Offset)
	fmt.Fprintf(os.Stderr, "!!! After log.Printf !!!\n")

	// Parse shop ID (empty string means all shops)
	var shopID uuid.UUID
	if req.ShopId != "" {
		var err error
		shopID, err = uuid.Parse(req.ShopId)
		if err != nil {
			log.Printf("Invalid shop_id: %v", err)
			return nil, err
		}
	}

	// Execute usecase
	output, err := h.list_productsUsecase.Execute(ctx, usecase.ListProductsInput{
		ShopID:        shopID,
		Category:      req.Category,
		PublishedOnly: req.PublishedOnly,
		Limit:         int(req.Limit),
		Offset:        int(req.Offset),
	})
	if err != nil {
		log.Printf("ListProducts usecase error: %v", err)
		return nil, err
	}

	// Convert to proto
	pbProducts := make([]*pb.Product, 0, len(output.Products))
	for _, product := range output.Products {
		pbProduct := &pb.Product{
			Id:            product.Id.String(),
			ShopId:        product.ShopId.String(),
			Name:          product.Name,
			Description:   product.Description,
			Price:         product.Price.String(),
			Category:      product.Category,
			StockQuantity: int32(product.StockQuantity),
			Published:     product.Published,
		}

		if product.Weight != nil {
			pbProduct.Weight = product.Weight.String()
		}
		if product.Size != nil {
			pbProduct.Size = *product.Size
		}
		if product.JanCode != nil {
			pbProduct.JanCode = *product.JanCode
		}

		pbProducts = append(pbProducts, pbProduct)
	}

	log.Printf("ListProducts returning %d products", len(pbProducts))
	return &pb.ListProductsResponse{
		Products: pbProducts,
	}, nil
}

// UploadProductImage handles UploadProductImage RPC
func (h *ShopServiceHandler) UploadProductImage(
	ctx context.Context,
	req *pb.UploadProductImageRequest,
) (*pb.UploadProductImageResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.UploadProductImageResponse{}, nil
}

// ManageVariation handles ManageVariation RPC
func (h *ShopServiceHandler) ManageVariation(
	ctx context.Context,
	req *pb.ManageVariationRequest,
) (*pb.ManageVariationResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.ManageVariationResponse{}, nil
}

// ListOrders handles ListOrders RPC
func (h *ShopServiceHandler) ListOrders(
	ctx context.Context,
	req *pb.ListOrdersRequest,
) (*pb.ListOrdersResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.ListOrdersResponse{}, nil
}

// GetOrderDetail handles GetOrderDetail RPC
func (h *ShopServiceHandler) GetOrderDetail(
	ctx context.Context,
	req *pb.GetOrderDetailRequest,
) (*pb.GetOrderDetailResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.GetOrderDetailResponse{}, nil
}

// UpdateOrderStatus handles UpdateOrderStatus RPC
func (h *ShopServiceHandler) UpdateOrderStatus(
	ctx context.Context,
	req *pb.UpdateOrderStatusRequest,
) (*pb.UpdateOrderStatusResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.UpdateOrderStatusResponse{}, nil
}

// GetSalesReport handles GetSalesReport RPC
func (h *ShopServiceHandler) GetSalesReport(
	ctx context.Context,
	req *pb.GetSalesReportRequest,
) (*pb.GetSalesReportResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.GetSalesReportResponse{}, nil
}

// ExportSalesData handles ExportSalesData RPC
func (h *ShopServiceHandler) ExportSalesData(
	ctx context.Context,
	req *pb.ExportSalesDataRequest,
) (*pb.ExportSalesDataResponse, error) {
	// TODO: Implement handler logic
	// 1. Convert request to usecase input
	// 2. Execute usecase
	// 3. Convert usecase output to response

	return &pb.ExportSalesDataResponse{}, nil
}

