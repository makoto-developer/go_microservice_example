package grpc

import (
	"context"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/repository"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/payment/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateCODPayment は代金引換の決済(支払い待ち)を作成する。
func (h *PaymentServiceHandler) CreateCODPayment(ctx context.Context, req *pb.CreateCODPaymentRequest) (*pb.CreateCODPaymentResponse, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}
	amount, err := strconv.Atoi(req.GetAmount())
	if err != nil || amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid amount")
	}

	output, err := h.codPaymentUsecase.Create(ctx, usecase.CreateCODPaymentInput{
		OrderID: orderID,
		Amount:  amount,
	})
	if err != nil {
		return nil, domainError(err)
	}

	return &pb.CreateCODPaymentResponse{
		PaymentId: output.PaymentID.String(),
		Message:   "COD payment created. It will be settled on delivery.",
	}, nil
}

// ConfirmCODPayment は配達完了時の集金を入金として確定する。
func (h *PaymentServiceHandler) ConfirmCODPayment(ctx context.Context, req *pb.ConfirmCODPaymentRequest) (*pb.ConfirmCODPaymentResponse, error) {
	paymentID, err := uuid.Parse(req.GetPaymentId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payment ID")
	}
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	if err := h.codPaymentUsecase.Confirm(ctx, usecase.ConfirmCODPaymentInput{
		PaymentID: paymentID,
		OrderID:   orderID,
	}); err != nil {
		return nil, domainError(err)
	}

	return &pb.ConfirmCODPaymentResponse{
		Success: true,
		Message: "COD payment confirmed",
	}, nil
}

// CreateRefund は決済に対する返金を行う。payment_id か order_id のどちらかで決済を特定できる。
func (h *PaymentServiceHandler) CreateRefund(ctx context.Context, req *pb.CreateRefundRequest) (*pb.CreateRefundResponse, error) {
	input := usecase.RefundPaymentInput{Reason: req.GetReason()}

	if req.GetPaymentId() != "" {
		id, err := uuid.Parse(req.GetPaymentId())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid payment ID")
		}
		input.PaymentID = id
	}
	if req.GetOrderId() != "" {
		id, err := uuid.Parse(req.GetOrderId())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid order ID")
		}
		input.OrderID = id
	}
	if input.PaymentID == uuid.Nil && input.OrderID == uuid.Nil {
		return nil, status.Error(codes.InvalidArgument, "payment_id or order_id is required")
	}
	if req.GetAmount() != "" {
		amount, err := strconv.Atoi(req.GetAmount())
		if err != nil || amount < 0 {
			return nil, status.Error(codes.InvalidArgument, "invalid amount")
		}
		input.Amount = amount // 0 = 全額返金
	}

	output, err := h.refundPaymentUsecase.Execute(ctx, input)
	if err != nil {
		return nil, domainError(err)
	}

	return &pb.CreateRefundResponse{
		RefundId: output.RefundID.String(),
		Message:  "Refund processed",
	}, nil
}

// GetRefundStatus は返金の状態を返す。
func (h *PaymentServiceHandler) GetRefundStatus(ctx context.Context, req *pb.GetRefundStatusRequest) (*pb.GetRefundStatusResponse, error) {
	refundID, err := uuid.Parse(req.GetRefundId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid refund ID")
	}

	refund, err := h.refundRepo.GetByID(ctx, refundID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if refund == nil {
		return nil, status.Error(codes.NotFound, "refund not found")
	}

	return &pb.GetRefundStatusResponse{
		Status:       refundStatusToProto(refund.Status),
		RefundAmount: int64(refund.Amount),
	}, nil
}

// ListPayments は決済の一覧を返す(管理者・加盟店画面用)。
// customer_id / shop_id / 日付フィルタはこのサンプルでは未対応(payments テーブルに列がない)。
func (h *PaymentServiceHandler) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	filter := repository.PaymentListFilter{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
	}
	if req.GetOrderId() != "" {
		id, err := uuid.Parse(req.GetOrderId())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid order ID")
		}
		filter.OrderID = id
	}
	for _, s := range req.GetStatusFilter() {
		if ds, ok := protoToDomainStatus(s); ok {
			filter.Statuses = append(filter.Statuses, ds)
		}
	}
	if m, ok := protoToDomainMethod(req.GetPaymentMethod()); ok {
		filter.Method = m
	}

	payments, total, err := h.paymentRepo.List(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	items := make([]*pb.Payment, 0, len(payments))
	for _, p := range payments {
		items = append(items, paymentToProto(p))
	}
	return &pb.ListPaymentsResponse{Payments: items, TotalCount: int32(total)}, nil
}

// GetPaymentDetail は決済の詳細を返す。
// requester_id / requester_role による権限チェックはサンプルのため省略。
func (h *PaymentServiceHandler) GetPaymentDetail(ctx context.Context, req *pb.GetPaymentDetailRequest) (*pb.GetPaymentDetailResponse, error) {
	paymentID, err := uuid.Parse(req.GetPaymentId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payment ID")
	}

	payment, err := h.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if payment == nil {
		return nil, status.Error(codes.NotFound, "payment not found")
	}

	return &pb.GetPaymentDetailResponse{Payment: paymentToProto(payment)}, nil
}

// --- 変換ヘルパー ---

func paymentToProto(p *domain.Payment) *pb.Payment {
	return &pb.Payment{
		Id:            p.ID.String(),
		OrderId:       p.OrderID.String(),
		PaymentMethod: methodToProto(p.PaymentMethod),
		Amount:        strconv.Itoa(p.Amount),
		Currency:      "jpy",
		Status:        statusToProto(p.Status),
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     timestamppb.New(p.UpdatedAt),
	}
}

func statusToProto(s domain.PaymentStatus) pb.PaymentStatus {
	switch s {
	case domain.PaymentStatusPending:
		return pb.PaymentStatus_PAYMENT_STATUS_PENDING
	case domain.PaymentStatusCompleted:
		return pb.PaymentStatus_PAYMENT_STATUS_SUCCEEDED
	case domain.PaymentStatusFailed:
		return pb.PaymentStatus_PAYMENT_STATUS_FAILED
	case domain.PaymentStatusRefunded:
		return pb.PaymentStatus_PAYMENT_STATUS_REFUNDED
	default:
		return pb.PaymentStatus_PAYMENT_STATUS_PAYMENT_STATUS_UNSPECIFIED
	}
}

func protoToDomainStatus(s pb.PaymentStatus) (domain.PaymentStatus, bool) {
	switch s {
	case pb.PaymentStatus_PAYMENT_STATUS_PENDING:
		return domain.PaymentStatusPending, true
	case pb.PaymentStatus_PAYMENT_STATUS_SUCCEEDED:
		return domain.PaymentStatusCompleted, true
	case pb.PaymentStatus_PAYMENT_STATUS_FAILED:
		return domain.PaymentStatusFailed, true
	case pb.PaymentStatus_PAYMENT_STATUS_REFUNDED:
		return domain.PaymentStatusRefunded, true
	default:
		return "", false
	}
}

func methodToProto(m domain.PaymentMethod) pb.PaymentMethodType {
	switch m {
	case domain.PaymentMethodCreditCard, domain.PaymentMethodDebitCard:
		return pb.PaymentMethodType_CREDIT_CARD
	case domain.PaymentMethodCashOnDelivery:
		return pb.PaymentMethodType_CASH_ON_DELIVERY
	default:
		return pb.PaymentMethodType_PAYMENT_METHOD_TYPE_UNSPECIFIED
	}
}

func protoToDomainMethod(m pb.PaymentMethodType) (domain.PaymentMethod, bool) {
	switch m {
	case pb.PaymentMethodType_CREDIT_CARD:
		return domain.PaymentMethodCreditCard, true
	case pb.PaymentMethodType_CASH_ON_DELIVERY:
		return domain.PaymentMethodCashOnDelivery, true
	default:
		return "", false
	}
}

func refundStatusToProto(s domain.RefundStatus) pb.RefundStatus {
	switch s {
	case domain.RefundStatusPending:
		return pb.RefundStatus_REFUND_STATUS_PENDING
	case domain.RefundStatusSucceeded:
		return pb.RefundStatus_REFUND_STATUS_SUCCEEDED
	case domain.RefundStatusFailed:
		return pb.RefundStatus_REFUND_STATUS_FAILED
	default:
		return pb.RefundStatus_REFUND_STATUS_REFUND_STATUS_UNSPECIFIED
	}
}

// domainError はドメインエラーを gRPC のステータスコードにマップする。
func domainError(err error) error {
	switch {
	case errors.Is(err, domain.ErrPaymentNotFound), errors.Is(err, domain.ErrRefundNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidAmount), errors.Is(err, domain.ErrInvalidRefundAmount):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrPaymentAlreadyProcessed):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrRefundNotAllowed), errors.Is(err, domain.ErrOrderMismatch), errors.Is(err, domain.ErrNotCODPayment):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
