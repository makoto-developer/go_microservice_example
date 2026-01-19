package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/repository"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/proto/order_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServiceHandler struct {
	pb.UnimplementedOrderServiceServer
	createOrderUsecase usecase.CreateOrderUsecase
	cancelOrderUsecase usecase.CancelOrderUsecase
	orderRepo          repository.OrderRepository
}

func NewOrderServiceHandler(
	createOrderUsecase usecase.CreateOrderUsecase,
	cancelOrderUsecase usecase.CancelOrderUsecase,
	orderRepo repository.OrderRepository,
) *OrderServiceHandler {
	return &OrderServiceHandler{
		createOrderUsecase: createOrderUsecase,
		cancelOrderUsecase: cancelOrderUsecase,
		orderRepo:          orderRepo,
	}
}

func (h *OrderServiceHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	// Convert cart items
	items := make([]usecase.OrderItemInput, len(req.CartItems))
	for i, item := range req.CartItems {
		productID, err := uuid.Parse(item.ProductId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid product ID")
		}

		items[i] = usecase.OrderItemInput{
			ProductID: productID,
			ShopID:    uuid.New(), // TODO: Get from product or cart item
			Quantity:  int(item.Quantity),
			Price:     int(item.UnitPrice),
		}
	}

	input := usecase.CreateOrderInput{
		CustomerID: customerID,
		Items:      items,
	}

	output, err := h.createOrderUsecase.Execute(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateOrderResponse{
		OrderId:     output.OrderID.String(),
		OrderNumber: output.OrderNumber,
		Message:     "Order created successfully",
	}, nil
}

func (h *OrderServiceHandler) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	input := usecase.CancelOrderInput{
		OrderID: orderID,
	}

	output, err := h.cancelOrderUsecase.Execute(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CancelOrderResponse{
		Success: output.Cancelled,
		Message: "Order cancelled successfully",
	}, nil
}
