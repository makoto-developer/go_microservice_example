package grpc

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/makoto-developer/go_microservice_example/proto/inventory-service/v1"
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
	inventoryID, err := uuid.Parse(req.InventoryId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid inventory_id: %v", err)
	}

	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}

	// For now, use the inventory_id directly
	// In a real implementation, you'd fetch the inventory to get product/variation IDs
	input := usecase.ReserveInventoryInput{
		ProductID:   inventoryID, // Simplified for now
		VariationID: nil,
		OrderID:     orderID,
		Quantity:    int(req.Quantity),
	}

	_, err = h.inventoryMgmt.ReserveInventory(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reserve: %v", err)
	}

	return &pb.ReserveStockResponse{}, nil
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

// Stub implementations for remaining RPCs
func (h *InventoryServiceHandler) RegisterInventory(ctx context.Context, req *pb.RegisterInventoryRequest) (*pb.RegisterInventoryResponse, error) {
	return &pb.RegisterInventoryResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (h *InventoryServiceHandler) UpdateInventoryQuantity(ctx context.Context, req *pb.UpdateInventoryQuantityRequest) (*pb.UpdateInventoryQuantityResponse, error) {
	return &pb.UpdateInventoryQuantityResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (h *InventoryServiceHandler) GetInventory(ctx context.Context, req *pb.GetInventoryRequest) (*pb.GetInventoryResponse, error) {
	return &pb.GetInventoryResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (h *InventoryServiceHandler) GetInventoryByProduct(ctx context.Context, req *pb.GetInventoryByProductRequest) (*pb.GetInventoryByProductResponse, error) {
	return &pb.GetInventoryByProductResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (h *InventoryServiceHandler) BulkGetInventory(ctx context.Context, req *pb.BulkGetInventoryRequest) (*pb.BulkGetInventoryResponse, error) {
	return &pb.BulkGetInventoryResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (h *InventoryServiceHandler) BulkReserveStock(ctx context.Context, req *pb.BulkReserveStockRequest) (*pb.BulkReserveStockResponse, error) {
	return &pb.BulkReserveStockResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (h *InventoryServiceHandler) ConfirmStock(ctx context.Context, req *pb.ConfirmStockRequest) (*pb.ConfirmStockResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}

	if err := h.inventoryMgmt.ConfirmInventory(ctx, orderID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to confirm: %v", err)
	}

	return &pb.ConfirmStockResponse{}, nil
}

func (h *InventoryServiceHandler) ReleaseExpiredReservations(ctx context.Context, req *pb.ReleaseExpiredReservationsRequest) (*pb.ReleaseExpiredReservationsResponse, error) {
	return &pb.ReleaseExpiredReservationsResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (h *InventoryServiceHandler) CheckStockAlert(ctx context.Context, req *pb.CheckStockAlertRequest) (*pb.CheckStockAlertResponse, error) {
	return &pb.CheckStockAlertResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (h *InventoryServiceHandler) GetInventoryHistory(ctx context.Context, req *pb.GetInventoryHistoryRequest) (*pb.GetInventoryHistoryResponse, error) {
	return &pb.GetInventoryHistoryResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (h *InventoryServiceHandler) RecordStockTaking(ctx context.Context, req *pb.RecordStockTakingRequest) (*pb.RecordStockTakingResponse, error) {
	return &pb.RecordStockTakingResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (h *InventoryServiceHandler) GetStockTakingHistory(ctx context.Context, req *pb.GetStockTakingHistoryRequest) (*pb.GetStockTakingHistoryResponse, error) {
	return &pb.GetStockTakingHistoryResponse{}, status.Error(codes.Unimplemented, "not implemented")
}
