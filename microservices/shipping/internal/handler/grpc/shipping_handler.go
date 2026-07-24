package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/client"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/repository"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/shipping/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ShippingServiceHandler struct {
	pb.UnimplementedShippingServiceServer
	createShipmentUsecase usecase.CreateShipmentUsecase
	shipmentRepo          repository.ShipmentRepository
	paymentClient         client.PaymentClient      // nil の場合は決済連携をスキップ
	notificationClient    client.NotificationClient // nil の場合は通知をスキップ
}

func NewShippingServiceHandler(
	createShipmentUsecase usecase.CreateShipmentUsecase,
	shipmentRepo repository.ShipmentRepository,
	paymentClient client.PaymentClient,
) *ShippingServiceHandler {
	return &ShippingServiceHandler{
		createShipmentUsecase: createShipmentUsecase,
		shipmentRepo:          shipmentRepo,
		paymentClient:         paymentClient,
	}
}

// WithNotification は配達完了通知(notification サービス連携)を有効にする。
func (h *ShippingServiceHandler) WithNotification(nc client.NotificationClient) *ShippingServiceHandler {
	h.notificationClient = nc
	return h
}

// CalculateShippingFee は配送料を見積もる(サンプルのため簡易テーブル)。
func (h *ShippingServiceHandler) CalculateShippingFee(ctx context.Context, req *pb.CalculateShippingFeeRequest) (*pb.CalculateShippingFeeResponse, error) {
	fee := 500
	if req.GetShippingMethod() == "express" {
		fee = 1000
	}
	if req.GetTotalWeightGrams() > 10_000 {
		fee += 500 // 10kg 超は大型扱い
	}

	return &pb.CalculateShippingFeeResponse{
		Success: true,
		Message: "Shipping fee calculated successfully",
		Fee:     fmt.Sprintf("%d", fee),
	}, nil
}

func (h *ShippingServiceHandler) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}
	customerID, err := uuid.Parse(req.GetCustomerId())
	if err != nil {
		customerID = uuid.Nil // 顧客 ID 無しでも出荷は作れる(バックオフィス起票)
	}
	carrier := req.GetCarrier()
	if carrier == "" {
		carrier = "yamato"
	}

	input := usecase.CreateShipmentInput{
		OrderID:         orderID,
		CustomerID:      customerID,
		ShippingAddress: req.GetShippingAddress(),
		Carrier:         carrier,
	}

	output, err := h.createShipmentUsecase.Execute(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateShipmentResponse{
		Success:    true,
		Message:    "Shipment created",
		ShipmentId: output.ShipmentID.String(),
	}, nil
}

// RegisterTrackingNumber は追跡番号を登録し、出荷を shipped にする。
func (h *ShippingServiceHandler) RegisterTrackingNumber(ctx context.Context, req *pb.RegisterTrackingNumberRequest) (*pb.RegisterTrackingNumberResponse, error) {
	if req.GetTrackingNumber() == "" {
		return nil, status.Error(codes.InvalidArgument, "tracking number is required")
	}
	shipment, err := h.getShipment(ctx, req.GetShipmentId())
	if err != nil {
		return nil, err
	}

	if err := h.shipmentRepo.UpdateTracking(ctx, shipment.ID, req.GetTrackingNumber()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := h.shipmentRepo.UpdateStatus(ctx, shipment.ID, domain.ShipmentStatusShipped); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterTrackingNumberResponse{
		Success: true,
		Message: "Tracking number registered",
	}, nil
}

// UpdateShipmentStatus は配送業者からの状態通知を反映する。
// delivered になったタイミングで、代引き注文なら決済サービスに集金確定を通知する。
func (h *ShippingServiceHandler) UpdateShipmentStatus(ctx context.Context, req *pb.UpdateShipmentStatusRequest) (*pb.UpdateShipmentStatusResponse, error) {
	shipment, err := h.getShipment(ctx, req.GetShipmentId())
	if err != nil {
		return nil, err
	}
	if !shipment.CanUpdate() {
		return nil, status.Error(codes.FailedPrecondition, "shipment can no longer be updated")
	}

	newStatus, ok := statusFromProto(req.GetNewStatus())
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid shipment status")
	}

	if err := h.shipmentRepo.UpdateStatus(ctx, shipment.ID, newStatus); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	message := "Shipment status updated"
	if newStatus == domain.ShipmentStatusDelivered {
		if h.paymentClient != nil {
			// 配達完了 = 代引きの集金完了。決済サービスに入金を確定させる
			if err := h.paymentClient.ConfirmCODByOrder(ctx, shipment.OrderID.String()); err != nil {
				// 集金確定の失敗で配送状態の更新は巻き戻さない(再送はバッチ想定)
				log.Printf("COD confirmation failed for order %s: %v", shipment.OrderID, err)
				message = "Shipment delivered, but COD confirmation failed (will retry)"
			} else {
				message = "Shipment delivered. COD payment (if any) confirmed."
			}
		}
		// 配達完了メール(best effort)
		if h.notificationClient != nil {
			if err := h.notificationClient.NotifyDelivered(
				ctx, shipment.CustomerID.String(), shipment.OrderID.String(), shipment.TrackingNumber,
			); err != nil {
				log.Printf("delivery notification failed for order %s: %v", shipment.OrderID, err)
			}
		}
	}

	return &pb.UpdateShipmentStatusResponse{
		Success: true,
		Message: message,
	}, nil
}

func (h *ShippingServiceHandler) GetShipmentDetail(ctx context.Context, req *pb.GetShipmentDetailRequest) (*pb.GetShipmentDetailResponse, error) {
	shipment, err := h.getShipment(ctx, req.GetShipmentId())
	if err != nil {
		return nil, err
	}

	return &pb.GetShipmentDetailResponse{
		Success:  true,
		Message:  "Shipment found",
		Shipment: shipmentToProto(shipment),
	}, nil
}

func (h *ShippingServiceHandler) GetShipmentByOrder(ctx context.Context, req *pb.GetShipmentByOrderRequest) (*pb.GetShipmentByOrderResponse, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
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
		Success:  true,
		Message:  "Shipment found",
		Shipment: shipmentToProto(shipment),
	}, nil
}

// SyncCarrierTracking は配送業者 API から追跡情報を取り込む(サンプルのためモック)。
func (h *ShippingServiceHandler) SyncCarrierTracking(ctx context.Context, req *pb.SyncCarrierTrackingRequest) (*pb.SyncCarrierTrackingResponse, error) {
	if _, err := h.getShipment(ctx, req.GetShipmentId()); err != nil {
		return nil, err
	}
	return &pb.SyncCarrierTrackingResponse{
		Success: true,
		Message: "Carrier tracking synced",
	}, nil
}

func (h *ShippingServiceHandler) ValidateAddress(ctx context.Context, req *pb.ValidateAddressRequest) (*pb.ValidateAddressResponse, error) {
	if req.GetPostalCode() == "" || req.GetPrefecture() == "" {
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
	normalized := fmt.Sprintf("〒%s %s%s%s",
		req.GetPostalCode(), req.GetPrefecture(), req.GetCity(), req.GetAddressLine())

	return &pb.NormalizeAddressResponse{
		Success:           true,
		Message:           "Address normalized",
		NormalizedAddress: normalized,
	}, nil
}

// --- ヘルパー ---

func (h *ShippingServiceHandler) getShipment(ctx context.Context, id string) (*domain.Shipment, error) {
	shipmentID, err := uuid.Parse(id)
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
	return shipment, nil
}

func statusFromProto(s pb.ShipmentStatus) (domain.ShipmentStatus, bool) {
	switch s {
	case pb.ShipmentStatus_SHIPMENT_STATUS_PENDING:
		return domain.ShipmentStatusPending, true
	case pb.ShipmentStatus_SHIPMENT_STATUS_PREPARING:
		return domain.ShipmentStatusPreparing, true
	case pb.ShipmentStatus_SHIPMENT_STATUS_SHIPPED:
		return domain.ShipmentStatusShipped, true
	case pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT:
		return domain.ShipmentStatusInTransit, true
	case pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED:
		return domain.ShipmentStatusDelivered, true
	case pb.ShipmentStatus_SHIPMENT_STATUS_FAILED:
		return domain.ShipmentStatusFailed, true
	default:
		return "", false
	}
}

func statusToProto(s domain.ShipmentStatus) pb.ShipmentStatus {
	switch s {
	case domain.ShipmentStatusPending:
		return pb.ShipmentStatus_SHIPMENT_STATUS_PENDING
	case domain.ShipmentStatusPreparing:
		return pb.ShipmentStatus_SHIPMENT_STATUS_PREPARING
	case domain.ShipmentStatusShipped:
		return pb.ShipmentStatus_SHIPMENT_STATUS_SHIPPED
	case domain.ShipmentStatusInTransit:
		return pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT
	case domain.ShipmentStatusDelivered:
		return pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED
	case domain.ShipmentStatusFailed:
		return pb.ShipmentStatus_SHIPMENT_STATUS_FAILED
	default:
		return pb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED
	}
}

func shipmentToProto(s *domain.Shipment) *pb.Shipment {
	out := &pb.Shipment{
		Id:                s.ID.String(),
		OrderId:           s.OrderID.String(),
		CustomerId:        s.CustomerID.String(),
		Status:            statusToProto(s.Status),
		TrackingNumber:    s.TrackingNumber,
		Carrier:           s.Carrier,
		ShippingAddress:   s.ShippingAddress,
		EstimatedDelivery: timestamppb.New(s.EstimatedDelivery),
		CreatedAt:         timestamppb.New(s.CreatedAt),
		UpdatedAt:         timestamppb.New(s.UpdatedAt),
	}
	if s.ActualDelivery != nil {
		out.ActualDelivery = timestamppb.New(*s.ActualDelivery)
	}
	return out
}
