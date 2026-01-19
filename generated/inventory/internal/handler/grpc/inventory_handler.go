package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/inventory/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/proto/inventory_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryServiceHandler struct {
	pb.UnimplementedInventoryServiceServer
	reserveStockUsecase usecase.ReserveStockUsecase
	releaseStockUsecase usecase.ReleaseStockUsecase
}

func NewInventoryServiceHandler(
	reserveStockUsecase usecase.ReserveStockUsecase,
	releaseStockUsecase usecase.ReleaseStockUsecase,
) *InventoryServiceHandler {
	return &InventoryServiceHandler{
		reserveStockUsecase: reserveStockUsecase,
		releaseStockUsecase: releaseStockUsecase,
	}
}

func (h *InventoryServiceHandler) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
	inventoryID, err := uuid.Parse(req.InventoryId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid inventory ID")
	}

	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	input := usecase.ReserveStockInput{
		ProductID: inventoryID,
		ShopID:    uuid.New(),
		Quantity:  int(req.Quantity),
		OrderID:   orderID,
	}

	_, err = h.reserveStockUsecase.Execute(ctx, input)
	if err != nil {
		if err == usecase.ErrInsufficientStock {
			return nil, status.Error(codes.FailedPrecondition, "insufficient stock")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ReserveStockResponse{
		Message: "Stock reserved successfully",
	}, nil
}

func (h *InventoryServiceHandler) ReleaseStock(ctx context.Context, req *pb.ReleaseStockRequest) (*pb.ReleaseStockResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	input := usecase.ReleaseStockInput{
		OrderID: orderID,
	}

	_, err = h.releaseStockUsecase.Execute(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ReleaseStockResponse{
		Success: true,
		Message: "Stock released successfully",
	}, nil
}
