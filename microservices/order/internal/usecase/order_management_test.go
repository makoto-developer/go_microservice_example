package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/client"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/domain"
)

// ---- フェイク ----

type fakeOrderRepo struct {
	created   []*domain.Order
	statuses  map[uuid.UUID][]domain.OrderStatus
	statusErr error
}

func newFakeOrderRepo() *fakeOrderRepo {
	return &fakeOrderRepo{statuses: map[uuid.UUID][]domain.OrderStatus{}}
}

func (r *fakeOrderRepo) Create(_ context.Context, order *domain.Order) error {
	r.created = append(r.created, order)
	return nil
}

func (r *fakeOrderRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Order, error) {
	for _, o := range r.created {
		if o.ID == id {
			return o, nil
		}
	}
	return nil, errors.New("not found")
}

func (r *fakeOrderRepo) GetByCustomerID(_ context.Context, _ uuid.UUID) ([]*domain.Order, error) {
	return r.created, nil
}

func (r *fakeOrderRepo) UpdateStatus(_ context.Context, id uuid.UUID, status domain.OrderStatus) error {
	if r.statusErr != nil {
		return r.statusErr
	}
	r.statuses[id] = append(r.statuses[id], status)
	return nil
}

func (r *fakeOrderRepo) lastStatus(id uuid.UUID) (domain.OrderStatus, bool) {
	history := r.statuses[id]
	if len(history) == 0 {
		return "", false
	}
	return history[len(history)-1], true
}

type fakeOrderItemRepo struct {
	items []*domain.OrderItem
}

func (r *fakeOrderItemRepo) Create(_ context.Context, item *domain.OrderItem) error {
	r.items = append(r.items, item)
	return nil
}

func (r *fakeOrderItemRepo) GetByOrderID(_ context.Context, _ uuid.UUID) ([]*domain.OrderItem, error) {
	return r.items, nil
}

type refundCall struct {
	orderID string
	reason  string
}

type fakePaymentClient struct {
	calls       []client.ProcessPaymentInput
	codCalls    []client.CODPaymentInput
	refundCalls []refundCall
	err         error
	refundErr   error
}

func (c *fakePaymentClient) ProcessPayment(_ context.Context, in client.ProcessPaymentInput) (*client.PaymentResult, error) {
	c.calls = append(c.calls, in)
	if c.err != nil {
		return nil, c.err
	}
	return &client.PaymentResult{PaymentID: "pay_test"}, nil
}

func (c *fakePaymentClient) CreateCODPayment(_ context.Context, in client.CODPaymentInput) (*client.PaymentResult, error) {
	c.codCalls = append(c.codCalls, in)
	if c.err != nil {
		return nil, c.err
	}
	return &client.PaymentResult{PaymentID: "pay_cod_test"}, nil
}

func (c *fakePaymentClient) RefundByOrder(_ context.Context, orderID string, reason string) error {
	c.refundCalls = append(c.refundCalls, refundCall{orderID: orderID, reason: reason})
	return c.refundErr
}

func (c *fakePaymentClient) Close() error { return nil }

func testInput() CreateOrderInput {
	return CreateOrderInput{
		CustomerID:      uuid.New(),
		AddressID:       uuid.New(),
		ShippingFee:     500,
		PaymentMethodID: "pm_test_card",
		Items: []OrderItemInput{
			{ProductID: uuid.New(), Quantity: 2, UnitPrice: 1000},
			{ProductID: uuid.New(), Quantity: 1, UnitPrice: 500},
		},
	}
}

// ---- テスト ----

func TestCreateOrder_ProcessesPaymentAndMarksPaid(t *testing.T) {
	orderRepo := newFakeOrderRepo()
	itemRepo := &fakeOrderItemRepo{}
	payment := &fakePaymentClient{}
	u := NewOrderManagementUsecase(orderRepo, itemRepo, payment)

	input := testInput()
	orderID, err := u.CreateOrder(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if len(payment.calls) != 1 {
		t.Fatalf("expected 1 payment call, got %d", len(payment.calls))
	}
	call := payment.calls[0]
	// 合計 = 2×1000 + 1×500 + 送料500
	if call.Amount != 3000 {
		t.Errorf("expected payment amount 3000, got %d", call.Amount)
	}
	if call.OrderID != orderID.String() {
		t.Errorf("payment order id = %s, want %s", call.OrderID, orderID)
	}
	if call.Currency != "jpy" {
		t.Errorf("payment currency = %s, want jpy", call.Currency)
	}
	if call.PaymentMethodID != input.PaymentMethodID {
		t.Errorf("payment method id = %s, want %s", call.PaymentMethodID, input.PaymentMethodID)
	}

	if status, ok := orderRepo.lastStatus(orderID); !ok || status != domain.OrderStatusPaid {
		t.Errorf("order status = %v (ok=%v), want %v", status, ok, domain.OrderStatusPaid)
	}
	if len(itemRepo.items) != 2 {
		t.Errorf("expected 2 order items, got %d", len(itemRepo.items))
	}
}

func TestCreateOrder_PaymentFailureCancelsOrder(t *testing.T) {
	orderRepo := newFakeOrderRepo()
	itemRepo := &fakeOrderItemRepo{}
	payment := &fakePaymentClient{err: errors.New("card declined")}
	u := NewOrderManagementUsecase(orderRepo, itemRepo, payment)

	_, err := u.CreateOrder(context.Background(), testInput())
	if err == nil {
		t.Fatal("expected error when payment fails")
	}
	if !strings.Contains(err.Error(), "payment failed") {
		t.Errorf("error should mention payment failure: %v", err)
	}

	if len(orderRepo.created) != 1 {
		t.Fatalf("expected order to be created before payment, got %d", len(orderRepo.created))
	}
	orderID := orderRepo.created[0].ID
	if status, ok := orderRepo.lastStatus(orderID); !ok || status != domain.OrderStatusCancelled {
		t.Errorf("order status = %v (ok=%v), want %v", status, ok, domain.OrderStatusCancelled)
	}
}

func TestCreateOrder_CODCreatesPaymentAndConfirmsOrder(t *testing.T) {
	orderRepo := newFakeOrderRepo()
	itemRepo := &fakeOrderItemRepo{}
	payment := &fakePaymentClient{}
	u := NewOrderManagementUsecase(orderRepo, itemRepo, payment)

	input := testInput()
	input.PaymentMethod = PaymentMethodCashOnDelivery
	orderID, err := u.CreateOrder(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if len(payment.calls) != 0 {
		t.Errorf("card payment should not run for COD, got %d calls", len(payment.calls))
	}
	if len(payment.codCalls) != 1 {
		t.Fatalf("expected 1 COD payment call, got %d", len(payment.codCalls))
	}
	if payment.codCalls[0].Amount != 3000 {
		t.Errorf("COD amount = %d, want 3000", payment.codCalls[0].Amount)
	}

	// 代引きは配達時に支払うため Paid ではなく Confirmed
	if status, ok := orderRepo.lastStatus(orderID); !ok || status != domain.OrderStatusConfirmed {
		t.Errorf("order status = %v (ok=%v), want %v", status, ok, domain.OrderStatusConfirmed)
	}
}

func TestCancelOrder_RefundsPayment(t *testing.T) {
	orderRepo := newFakeOrderRepo()
	itemRepo := &fakeOrderItemRepo{}
	payment := &fakePaymentClient{}
	u := NewOrderManagementUsecase(orderRepo, itemRepo, payment)

	orderID := uuid.New()
	if err := u.CancelOrder(context.Background(), orderID); err != nil {
		t.Fatalf("CancelOrder returned error: %v", err)
	}

	if len(payment.refundCalls) != 1 {
		t.Fatalf("expected 1 refund call, got %d", len(payment.refundCalls))
	}
	if payment.refundCalls[0].orderID != orderID.String() {
		t.Errorf("refund order id = %s, want %s", payment.refundCalls[0].orderID, orderID)
	}
	if status, ok := orderRepo.lastStatus(orderID); !ok || status != domain.OrderStatusCancelled {
		t.Errorf("order status = %v (ok=%v), want %v", status, ok, domain.OrderStatusCancelled)
	}
}

func TestCancelOrder_RefundFailureKeepsOrder(t *testing.T) {
	orderRepo := newFakeOrderRepo()
	itemRepo := &fakeOrderItemRepo{}
	payment := &fakePaymentClient{refundErr: errors.New("payment service down")}
	u := NewOrderManagementUsecase(orderRepo, itemRepo, payment)

	orderID := uuid.New()
	err := u.CancelOrder(context.Background(), orderID)
	if err == nil {
		t.Fatal("expected error when refund fails")
	}
	// 返金に失敗したらキャンセルしない(返金漏れ防止)
	if _, ok := orderRepo.lastStatus(orderID); ok {
		t.Error("order should not be cancelled when refund fails")
	}
}

func TestCreateOrder_WithoutPaymentClientSkipsPayment(t *testing.T) {
	orderRepo := newFakeOrderRepo()
	itemRepo := &fakeOrderItemRepo{}
	u := NewOrderManagementUsecase(orderRepo, itemRepo, nil)

	orderID, err := u.CreateOrder(context.Background(), testInput())
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}
	if _, ok := orderRepo.lastStatus(orderID); ok {
		t.Error("status should not be updated when payment client is not configured")
	}
}
