package grpc

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/repository"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/payment/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// --- インメモリのフェイクリポジトリ ---

type fakePaymentRepo struct {
	payments map[uuid.UUID]*domain.Payment
}

func newFakePaymentRepo() *fakePaymentRepo {
	return &fakePaymentRepo{payments: map[uuid.UUID]*domain.Payment{}}
}

func (r *fakePaymentRepo) Create(_ context.Context, p *domain.Payment) error {
	cp := *p
	r.payments[p.ID] = &cp
	return nil
}

func (r *fakePaymentRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Payment, error) {
	if p, ok := r.payments[id]; ok {
		cp := *p
		return &cp, nil
	}
	return nil, nil
}

func (r *fakePaymentRepo) GetByOrderID(_ context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	for _, p := range r.payments {
		if p.OrderID == orderID {
			cp := *p
			return &cp, nil
		}
	}
	return nil, nil
}

func (r *fakePaymentRepo) List(_ context.Context, filter repository.PaymentListFilter) ([]*domain.Payment, int, error) {
	items := []*domain.Payment{}
	for _, p := range r.payments {
		if filter.OrderID != uuid.Nil && p.OrderID != filter.OrderID {
			continue
		}
		if filter.Method != "" && p.PaymentMethod != filter.Method {
			continue
		}
		cp := *p
		items = append(items, &cp)
	}
	return items, len(items), nil
}

func (r *fakePaymentRepo) UpdateStatus(_ context.Context, id uuid.UUID, s domain.PaymentStatus, txid string) error {
	if p, ok := r.payments[id]; ok {
		p.Status = s
		p.TransactionID = txid
	}
	return nil
}

type fakeRefundRepo struct {
	refunds map[uuid.UUID]*domain.Refund
}

func newFakeRefundRepo() *fakeRefundRepo {
	return &fakeRefundRepo{refunds: map[uuid.UUID]*domain.Refund{}}
}

func (r *fakeRefundRepo) Create(_ context.Context, refund *domain.Refund) error {
	cp := *refund
	r.refunds[refund.ID] = &cp
	return nil
}

func (r *fakeRefundRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Refund, error) {
	if refund, ok := r.refunds[id]; ok {
		cp := *refund
		return &cp, nil
	}
	return nil, nil
}

func newTestHandler() (*PaymentServiceHandler, *fakePaymentRepo, *fakeRefundRepo) {
	paymentRepo := newFakePaymentRepo()
	refundRepo := newFakeRefundRepo()
	h := NewPaymentServiceHandler(
		nil,
		usecase.NewCODPaymentUsecase(paymentRepo),
		usecase.NewRefundPaymentUsecase(paymentRepo, refundRepo),
		paymentRepo,
		refundRepo,
	)
	return h, paymentRepo, refundRepo
}

// --- COD ---

func TestCODPayment_FullFlow(t *testing.T) {
	h, repo, _ := newTestHandler()
	orderID := uuid.New()

	created, err := h.CreateCODPayment(context.Background(), &pb.CreateCODPaymentRequest{
		OrderId: orderID.String(),
		Amount:  "3200",
	})
	if err != nil {
		t.Fatalf("CreateCODPayment: %v", err)
	}

	paymentID := uuid.MustParse(created.GetPaymentId())
	if p := repo.payments[paymentID]; p == nil || p.Status != domain.PaymentStatusPending {
		t.Fatalf("COD payment should be pending until delivery, got %+v", repo.payments[paymentID])
	}

	// 配達完了 → 集金確定
	confirmed, err := h.ConfirmCODPayment(context.Background(), &pb.ConfirmCODPaymentRequest{
		PaymentId: created.GetPaymentId(),
		OrderId:   orderID.String(),
	})
	if err != nil {
		t.Fatalf("ConfirmCODPayment: %v", err)
	}
	if !confirmed.GetSuccess() {
		t.Error("expected success")
	}
	if p := repo.payments[paymentID]; p.Status != domain.PaymentStatusCompleted {
		t.Errorf("status = %s, want completed", p.Status)
	}
}

func TestCreateCODPayment_RejectsDuplicateOrder(t *testing.T) {
	h, _, _ := newTestHandler()
	orderID := uuid.New().String()

	if _, err := h.CreateCODPayment(context.Background(), &pb.CreateCODPaymentRequest{OrderId: orderID, Amount: "1000"}); err != nil {
		t.Fatalf("first CreateCODPayment: %v", err)
	}
	_, err := h.CreateCODPayment(context.Background(), &pb.CreateCODPaymentRequest{OrderId: orderID, Amount: "1000"})
	if status.Code(err) != codes.AlreadyExists {
		t.Errorf("expected AlreadyExists, got %v", err)
	}
}

func TestConfirmCODPayment_RejectsOrderMismatch(t *testing.T) {
	h, _, _ := newTestHandler()

	created, err := h.CreateCODPayment(context.Background(), &pb.CreateCODPaymentRequest{
		OrderId: uuid.New().String(),
		Amount:  "1000",
	})
	if err != nil {
		t.Fatalf("CreateCODPayment: %v", err)
	}

	_, err = h.ConfirmCODPayment(context.Background(), &pb.ConfirmCODPaymentRequest{
		PaymentId: created.GetPaymentId(),
		OrderId:   uuid.New().String(), // 別の注文
	})
	if status.Code(err) != codes.FailedPrecondition {
		t.Errorf("expected FailedPrecondition, got %v", err)
	}
}

func TestConfirmCODPayment_RejectsNonCOD(t *testing.T) {
	h, repo, _ := newTestHandler()

	p := domain.NewPayment(uuid.New(), 5000, domain.PaymentMethodCreditCard)
	repo.payments[p.ID] = p

	_, err := h.ConfirmCODPayment(context.Background(), &pb.ConfirmCODPaymentRequest{
		PaymentId: p.ID.String(),
		OrderId:   p.OrderID.String(),
	})
	if status.Code(err) != codes.FailedPrecondition {
		t.Errorf("expected FailedPrecondition, got %v", err)
	}
}

// --- 返金 ---

func TestCreateRefund_FullRefundByOrderID(t *testing.T) {
	h, repo, refundRepo := newTestHandler()

	p := domain.NewPayment(uuid.New(), 4800, domain.PaymentMethodCreditCard)
	if err := p.Complete("TXN-test"); err != nil {
		t.Fatal(err)
	}
	repo.payments[p.ID] = p

	resp, err := h.CreateRefund(context.Background(), &pb.CreateRefundRequest{
		OrderId: p.OrderID.String(),
		Reason:  "customer cancelled",
	})
	if err != nil {
		t.Fatalf("CreateRefund: %v", err)
	}

	refundID := uuid.MustParse(resp.GetRefundId())
	refund := refundRepo.refunds[refundID]
	if refund == nil {
		t.Fatal("refund not persisted")
	}
	if refund.Amount != 4800 {
		t.Errorf("refund amount = %d, want 4800 (全額返金)", refund.Amount)
	}
	if repo.payments[p.ID].Status != domain.PaymentStatusRefunded {
		t.Errorf("payment status = %s, want refunded", repo.payments[p.ID].Status)
	}

	// GetRefundStatus で参照できる
	st, err := h.GetRefundStatus(context.Background(), &pb.GetRefundStatusRequest{RefundId: resp.GetRefundId()})
	if err != nil {
		t.Fatalf("GetRefundStatus: %v", err)
	}
	if st.GetStatus() != pb.RefundStatus_REFUND_STATUS_SUCCEEDED {
		t.Errorf("refund status = %v, want SUCCEEDED", st.GetStatus())
	}
	if st.GetRefundAmount() != 4800 {
		t.Errorf("refund amount = %d, want 4800", st.GetRefundAmount())
	}
}

func TestCreateRefund_RejectsPendingPayment(t *testing.T) {
	h, repo, _ := newTestHandler()

	p := domain.NewPayment(uuid.New(), 1000, domain.PaymentMethodCashOnDelivery)
	repo.payments[p.ID] = p // pending のまま

	_, err := h.CreateRefund(context.Background(), &pb.CreateRefundRequest{
		PaymentId: p.ID.String(),
		Reason:    "test",
	})
	if status.Code(err) != codes.FailedPrecondition {
		t.Errorf("expected FailedPrecondition, got %v", err)
	}
}

func TestCreateRefund_RejectsOverAmount(t *testing.T) {
	h, repo, _ := newTestHandler()

	p := domain.NewPayment(uuid.New(), 1000, domain.PaymentMethodCreditCard)
	if err := p.Complete("TXN-test"); err != nil {
		t.Fatal(err)
	}
	repo.payments[p.ID] = p

	_, err := h.CreateRefund(context.Background(), &pb.CreateRefundRequest{
		PaymentId: p.ID.String(),
		Amount:    "1001",
		Reason:    "too much",
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", err)
	}
}

func TestCreateRefund_NotFound(t *testing.T) {
	h, _, _ := newTestHandler()

	_, err := h.CreateRefund(context.Background(), &pb.CreateRefundRequest{
		OrderId: uuid.New().String(),
		Reason:  "no payment",
	})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound, got %v", err)
	}
}

// --- 一覧・詳細 ---

func TestListPayments_FiltersByOrder(t *testing.T) {
	h, repo, _ := newTestHandler()

	target := domain.NewPayment(uuid.New(), 2000, domain.PaymentMethodCashOnDelivery)
	other := domain.NewPayment(uuid.New(), 9999, domain.PaymentMethodCreditCard)
	repo.payments[target.ID] = target
	repo.payments[other.ID] = other

	resp, err := h.ListPayments(context.Background(), &pb.ListPaymentsRequest{
		OrderId: target.OrderID.String(),
	})
	if err != nil {
		t.Fatalf("ListPayments: %v", err)
	}
	if len(resp.GetPayments()) != 1 {
		t.Fatalf("got %d payments, want 1", len(resp.GetPayments()))
	}
	got := resp.GetPayments()[0]
	if got.GetAmount() != strconv.Itoa(2000) {
		t.Errorf("amount = %s, want 2000", got.GetAmount())
	}
	if got.GetPaymentMethod() != pb.PaymentMethodType_CASH_ON_DELIVERY {
		t.Errorf("method = %v, want CASH_ON_DELIVERY", got.GetPaymentMethod())
	}
}

func TestGetPaymentDetail_ReturnsPayment(t *testing.T) {
	h, repo, _ := newTestHandler()

	p := domain.NewPayment(uuid.New(), 750, domain.PaymentMethodCashOnDelivery)
	repo.payments[p.ID] = p

	resp, err := h.GetPaymentDetail(context.Background(), &pb.GetPaymentDetailRequest{
		PaymentId: p.ID.String(),
	})
	if err != nil {
		t.Fatalf("GetPaymentDetail: %v", err)
	}
	pm := resp.GetPayment()
	if pm.GetId() != p.ID.String() || pm.GetStatus() != pb.PaymentStatus_PAYMENT_STATUS_PENDING {
		t.Errorf("unexpected payment: %+v", pm)
	}
	if pm.GetCurrency() != "jpy" {
		t.Errorf("currency = %s, want jpy", pm.GetCurrency())
	}
}
