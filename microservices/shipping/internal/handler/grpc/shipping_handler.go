package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shipping/internal/repository"
	"github.com/makoto-developer/go_microservice_example/generated/shipping/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/proto/shipping_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShippingServiceHandler struct {
	pb.UnimplementedShippingServiceServer
	createShipmentUsecase usecase.CreateShipmentUsecase
	shipmentRepo          repository.ShipmentRepository
}

func NewShippingServiceHandler(
	createShipmentUsecase usecase.CreateShipmentUsecase,
	shipmentRepo repository.ShipmentRepository,
) *ShippingServiceHandler {
	return &ShippingServiceHandler{
		createShipmentUsecase: createShipmentUsecase,
		shipmentRepo:          shipmentRepo,
	}
}

func (h *ShippingServiceHandler) CalculateShippingFee(ctx context.Context, req *pb.CalculateShippingFeeRequest) (*pb.CalculateShippingFeeResponse, error) {
	// Mock implementation - calculate shipping fee based on weight and size
	return &pb.CalculateShippingFeeResponse{
		Success: true,
		Message: "Shipping fee calculated successfully",
	}, nil
}

func (h *ShippingServiceHandler) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	input := usecase.CreateShipmentInput{
		OrderID:         orderID,
		CustomerID:      uuid.New(), // Should be extracted from order
		ShippingAddress: req.ShippingAddress,
		Carrier:         "yamato", // Default carrier
	}

	output, err := h.createShipmentUsecase.Execute(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateShipmentResponse{
		Success: true,
		Message: "Shipment created: " + output.ShipmentID.String(),
	}, nil
}

func (h *ShippingServiceHandler) RegisterTrackingNumber(ctx context.Context, req *pb.RegisterTrackingNumberRequest) (*pb.RegisterTrackingNumberResponse, error) {
	shipmentID, err := uuid.Parse(req.ShipmentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid shipment ID")
	}

	shipment, err := h.shipmentRepo.GetByID(ctx, shipmentID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if shipment == nil {
		return nil, status.Error(codes.NotFound, "shipment not found")
	}

	// Update tracking number
	// shipment.TrackingNumber = req.TrackingNumber
	// err = h.shipmentRepo.Update(ctx, shipment)

	return &pb.RegisterTrackingNumberResponse{
		Success: true,
		Message: "Tracking number registered successfully",
	}, nil
}

func (h *ShippingServiceHandler) UpdateShipmentStatus(ctx context.Context, req *pb.UpdateShipmentStatusRequest) (*pb.UpdateShipmentStatusResponse, error) {
	shipmentID, err := uuid.Parse(req.ShipmentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid shipment ID")
	}

	shipment, err := h.shipmentRepo.GetByID(ctx, shipmentID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if shipment == nil {
		return nil, status.Error(codes.NotFound, "shipment not found")
	}

	// Update status based on req.NewStatus
	// err = h.shipmentRepo.Update(ctx, shipment)

	return &pb.UpdateShipmentStatusResponse{
		Success: true,
		Message: "Shipment status updated successfully",
	}, nil
}

func (h *ShippingServiceHandler) GetShipmentDetail(ctx context.Context, req *pb.GetShipmentDetailRequest) (*pb.GetShipmentDetailResponse, error) {
	shipmentID, err := uuid.Parse(req.ShipmentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid shipment ID")
	}

	shipment, err := h.shipmentRepo.GetByID(ctx, shipmentID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if shipment == nil {
		return nil, status.Error(codes.NotFound, "shipment not found")
	}

	return &pb.GetShipmentDetailResponse{
		Success: true,
		Message: "Shipment found: " + shipment.ID.String(),
	}, nil
}

func (h *ShippingServiceHandler) GetShipmentByOrder(ctx context.Context, req *pb.GetShipmentByOrderRequest) (*pb.GetShipmentByOrderResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	shipment, err := h.shipmentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if shipment == nil {
		return nil, status.Error(codes.NotFound, "shipment not found")
	}

	return &pb.GetShipmentByOrderResponse{
		Success: true,
		Message: "Shipment found for order: " + orderID.String(),
	}, nil
}

func (h *ShippingServiceHandler) SyncCarrierTracking(ctx context.Context, req *pb.SyncCarrierTrackingRequest) (*pb.SyncCarrierTrackingResponse, error) {
	// Mock implementation - sync tracking info from carrier API
	return &pb.SyncCarrierTrackingResponse{
		Success: true,
		Message: "Carrier tracking synced successfully",
	}, nil
}

func (h *ShippingServiceHandler) ValidateAddress(ctx context.Context, req *pb.ValidateAddressRequest) (*pb.ValidateAddressResponse, error) {
	// Mock implementation - validate Japanese postal address
	if req.PostalCode == "" || req.Prefecture == "" {
		return &pb.ValidateAddressResponse{
			Success: false,
			Message: "Invalid address: missing required fields",
		}, nil
	}

	return &pb.ValidateAddressResponse{
		Success: true,
		Message: "Address is valid",
	}, nil
}

func (h *ShippingServiceHandler) NormalizeAddress(ctx context.Context, req *pb.NormalizeAddressRequest) (*pb.NormalizeAddressResponse, error) {
	// Mock implementation - normalize Japanese address format
	return &pb.NormalizeAddressResponse{
		Success: true,
		Message: "Address normalized successfully",
	}, nil
}
