package grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/order/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderServiceHandler struct {
	pb.UnimplementedOrderServiceServer
	orderMgmt usecase.OrderManagementUsecase
}

func NewOrderServiceHandler(orderMgmt usecase.OrderManagementUsecase) *OrderServiceHandler {
	return &OrderServiceHandler{orderMgmt: orderMgmt}
}

func (h *OrderServiceHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	customerID, err := uuid.Parse(req.GetCustomerId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid customer_id: %v", err)
	}

	addressID, err := uuid.Parse(req.GetShippingAddressId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid shipping_address_id: %v", err)
	}

	var items []usecase.OrderItemInput
	for _, item := range req.GetCartItems() {
		productID, err := uuid.Parse(item.GetProductId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid product_id: %v", err)
		}

		var variationID *uuid.UUID
		if item.GetVariationId() != "" {
			vid, err := uuid.Parse(item.GetVariationId())
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid variation_id: %v", err)
			}
			variationID = &vid
		}

		items = append(items, usecase.OrderItemInput{
			ProductID:   productID,
			VariationID: variationID,
			Quantity:    int(item.GetQuantity()),
			UnitPrice:   item.GetUnitPrice(),
		})
	}

	// Calculate shipping fee based on shipping method (simplified for now)
	shippingFee := int64(500) // Default shipping fee
	if req.GetShippingMethod() == "express" {
		shippingFee = 1000
	}

	// 支払い方法: 代引き指定のみ分岐し、未指定はクレジットカード扱い
	paymentMethod := usecase.PaymentMethodCreditCard
	if req.GetPaymentMethod() == pb.PaymentMethod_CASH_ON_DELIVERY {
		paymentMethod = usecase.PaymentMethodCashOnDelivery
	}

	input := usecase.CreateOrderInput{
		CustomerID:      customerID,
		AddressID:       addressID,
		Items:           items,
		ShippingFee:     shippingFee,
		PaymentMethod:   paymentMethod,
		PaymentMethodID: req.GetPaymentMethodId(),
	}

	orderID, err := h.orderMgmt.CreateOrder(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	message := "Order created and payment completed successfully"
	if paymentMethod == usecase.PaymentMethodCashOnDelivery {
		message = "Order created. Payment will be collected on delivery."
	}

	return &pb.CreateOrderResponse{
		OrderId: orderID.String(),
		Message: message,
	}, nil
}

func (h *OrderServiceHandler) GetOrderDetail(ctx context.Context, req *pb.GetOrderDetailRequest) (*pb.GetOrderDetailResponse, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}

	order, err := h.orderMgmt.GetOrder(ctx, orderID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}

	return &pb.GetOrderDetailResponse{Order: orderToProto(order)}, nil
}

// ListOrders は顧客の注文一覧を返す(注文履歴画面用)。
// shop_id 等でのサーバー側絞り込みはこのサンプルでは未対応。
func (h *OrderServiceHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	customerID, err := uuid.Parse(req.GetCustomerId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "customer_id is required: %v", err)
	}

	orders, err := h.orderMgmt.ListOrdersByCustomer(ctx, customerID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}

	items := make([]*pb.Order, 0, len(orders))
	for _, order := range orders {
		items = append(items, orderToProto(order))
	}

	return &pb.ListOrdersResponse{
		Orders:     items,
		TotalCount: int32(len(items)),
		Page:       1,
		PageSize:   int32(len(items)),
	}, nil
}

func (h *OrderServiceHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}

	domainStatus := convertProtoToOrderStatus(req.GetNewStatus())
	if err := h.orderMgmt.UpdateOrderStatus(ctx, orderID, domainStatus); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}

	return &pb.UpdateOrderStatusResponse{
		Success: true,
		Message: "Order status updated successfully",
	}, nil
}

func (h *OrderServiceHandler) GetOrderStatusHistory(ctx context.Context, req *pb.GetOrderStatusHistoryRequest) (*pb.GetOrderStatusHistoryResponse, error) {
	return &pb.GetOrderStatusHistoryResponse{
		History: []pb.OrderStatus{},
	}, nil
}

func (h *OrderServiceHandler) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}

	if err := h.orderMgmt.CancelOrder(ctx, orderID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to cancel order: %v", err)
	}

	return &pb.CancelOrderResponse{
		Success: true,
		Message: "Order cancelled successfully",
	}, nil
}

func (h *OrderServiceHandler) SearchOrders(ctx context.Context, req *pb.SearchOrdersRequest) (*pb.SearchOrdersResponse, error) {
	return &pb.SearchOrdersResponse{
		Orders:     []*pb.Order{},
		TotalCount: 0,
	}, nil
}

func (h *OrderServiceHandler) GetOrderStatistics(ctx context.Context, req *pb.GetOrderStatisticsRequest) (*pb.GetOrderStatisticsResponse, error) {
	return &pb.GetOrderStatisticsResponse{
		TotalOrders:     0,
		TotalSales:      0,
		PendingOrders:   0,
		CompletedOrders: 0,
	}, nil
}

func (h *OrderServiceHandler) GetProductSalesRanking(ctx context.Context, req *pb.GetProductSalesRankingRequest) (*pb.GetProductSalesRankingResponse, error) {
	return &pb.GetProductSalesRankingResponse{
		Rankings: []*pb.ProductSalesRank{},
	}, nil
}

func (h *OrderServiceHandler) ExportOrdersToCSV(ctx context.Context, req *pb.ExportOrdersToCSVRequest) (*pb.ExportOrdersToCSVResponse, error) {
	return &pb.ExportOrdersToCSVResponse{
		CsvUrl:  "",
		Message: "Export functionality not yet implemented",
	}, nil
}

func (h *OrderServiceHandler) CreateReorder(ctx context.Context, req *pb.CreateReorderRequest) (*pb.CreateReorderResponse, error) {
	return &pb.CreateReorderResponse{
		OrderId: "",
		Message: "Reorder functionality not yet implemented",
	}, nil
}

// Helper functions
func orderToProto(order *domain.Order) *pb.Order {
	// 注: domain.Order は支払い方法を保持しないため payment_method は未設定のまま。
	// 支払いの詳細は payment サービス側(ListPayments / GetPaymentDetail)が持つ
	return &pb.Order{
		Id:          order.ID.String(),
		OrderNumber: order.OrderNumber,
		CustomerId:  order.CustomerID.String(),
		Status:      convertOrderStatus(order.Status),
		TotalAmount: fmt.Sprintf("%d", order.TotalAmount),
		ShippingFee: fmt.Sprintf("%d", order.ShippingFee),
		CreatedAt:   timestamppb.New(order.CreatedAt),
		UpdatedAt:   timestamppb.New(order.UpdatedAt),
	}
}

func convertOrderStatus(status domain.OrderStatus) pb.OrderStatus {
	switch status {
	case domain.OrderStatusPending:
		return pb.OrderStatus_PENDING
	case domain.OrderStatusConfirmed:
		return pb.OrderStatus_CONFIRMED
	case domain.OrderStatusPaid:
		return pb.OrderStatus_PAYMENT_PROCESSING
	case domain.OrderStatusShipped:
		return pb.OrderStatus_SHIPPED
	case domain.OrderStatusDelivered:
		return pb.OrderStatus_DELIVERED
	case domain.OrderStatusCancelled:
		return pb.OrderStatus_CANCELLED
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

func convertProtoToOrderStatus(status pb.OrderStatus) domain.OrderStatus {
	switch status {
	case pb.OrderStatus_PENDING:
		return domain.OrderStatusPending
	case pb.OrderStatus_CONFIRMED:
		return domain.OrderStatusConfirmed
	case pb.OrderStatus_PAYMENT_PROCESSING:
		return domain.OrderStatusPaid
	case pb.OrderStatus_SHIPPED:
		return domain.OrderStatusShipped
	case pb.OrderStatus_DELIVERED:
		return domain.OrderStatusDelivered
	case pb.OrderStatus_CANCELLED:
		return domain.OrderStatusCancelled
	default:
		return domain.OrderStatusPending
	}
}
