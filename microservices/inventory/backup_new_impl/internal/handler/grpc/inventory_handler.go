package grpc

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/makoto-developer/go_microservice_example/proto/inventory_service/v1"
	"github.com/makoto-developer/go_microservice_example/generated/inventory/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryServiceHandler struct {
	pb.UnimplementedInventoryServiceServer
	inventoryMgmt usecase.InventoryManagementUsecase
}

func NewInventoryServiceHandler(inventoryMgmt usecase.InventoryManagementUsecase) *InventoryServiceHandler {
	return &InventoryServiceHandler{inventoryMgmt: inventoryMgmt}
}

func (h *InventoryServiceHandler) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product_id: %v", err)
	}
	
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}
	
	var variationID *uuid.UUID
	if req.VariationId != "" {
		vid, err := uuid.Parse(req.VariationId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid variation_id: %v", err)
		}
		variationID = &vid
	}
	
	input := usecase.ReserveInventoryInput{
		ProductID:   productID,
		VariationID: variationID,
		OrderID:     orderID,
		Quantity:    int(req.Quantity),
	}
	
	reservationID, err := h.inventoryMgmt.ReserveInventory(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reserve: %v", err)
	}
	
	return &pb.ReserveStockResponse{ReservationId: reservationID.String()}, nil
}

func (h *InventoryServiceHandler) ReleaseStock(ctx context.Context, req *pb.ReleaseStockRequest) (*pb.ReleaseStockResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}
	
	if err := h.inventoryMgmt.ReleaseInventory(ctx, orderID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to release: %v", err)
	}
	
	return &pb.ReleaseStockResponse{}, nil
}

// Stub implementations
func (h *InventoryServiceHandler) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
	return &pb.UpdateStockResponse{}, nil
}

func (h *InventoryServiceHandler) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.GetStockResponse, error) {
	return &pb.GetStockResponse{}, nil
}

func (h *InventoryServiceHandler) BulkReserveStock(ctx context.Context, req *pb.BulkReserveStockRequest) (*pb.BulkReserveStockResponse, error) {
	return &pb.BulkReserveStockResponse{}, nil
}

func (h *InventoryServiceHandler) ConfirmReservation(ctx context.Context, req *pb.ConfirmReservationRequest) (*pb.ConfirmReservationResponse, error) {
	return &pb.ConfirmReservationResponse{}, nil
}

func (h *InventoryServiceHandler) CheckStockAvailability(ctx context.Context, req *pb.CheckStockAvailabilityRequest) (*pb.CheckStockAvailabilityResponse, error) {
	return &pb.CheckStockAvailabilityResponse{}, nil
}

func (h *InventoryServiceHandler) GetLowStockAlerts(ctx context.Context, req *pb.GetLowStockAlertsRequest) (*pb.GetLowStockAlertsResponse, error) {
	return &pb.GetLowStockAlertsResponse{}, nil
}

func (h *InventoryServiceHandler) GetStockHistory(ctx context.Context, req *pb.GetStockHistoryRequest) (*pb.GetStockHistoryResponse, error) {
	return &pb.GetStockHistoryResponse{}, nil
}
