// サービス横断の統合テスト。
// payment / shipping / inventory のインメモリ gRPC サーバ(testsupport)を実際に起動し、
// order のユースケースから本物のクライアント・実 TCP 通信で一連のフローを検証する。
//
//	代引き購入   : order → inventory.BulkReserveStock → payment.CreateCODPayment(pending)
//	               → shipping.CreateShipment → inventory.ConfirmStock
//	配達完了     : shipping → payment.ListPayments → ConfirmCODPayment(completed)
//	カード購入   : order → payment.CreatePaymentIntent/ConfirmPayment(completed)
//	キャンセル   : order → payment.CreateRefund(refunded)
package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventorytest "github.com/makoto-developer/go_microservice_example/microservices/inventory/testsupport"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/client"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/repository"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/usecase"
	paymenttest "github.com/makoto-developer/go_microservice_example/microservices/payment/testsupport"
	shippingpb "github.com/makoto-developer/go_microservice_example/microservices/shipping/proto"
	shippingtest "github.com/makoto-developer/go_microservice_example/microservices/shipping/testsupport"
)

// ---- インメモリ order リポジトリ ----

type memOrderRepo struct {
	orders   map[uuid.UUID]*domain.Order
	statuses map[uuid.UUID]domain.OrderStatus
}

func newMemOrderRepo() *memOrderRepo {
	return &memOrderRepo{
		orders:   map[uuid.UUID]*domain.Order{},
		statuses: map[uuid.UUID]domain.OrderStatus{},
	}
}

func (r *memOrderRepo) Create(_ context.Context, o *domain.Order) error {
	r.orders[o.ID] = o
	return nil
}

func (r *memOrderRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Order, error) {
	return r.orders[id], nil
}

func (r *memOrderRepo) GetByCustomerID(_ context.Context, _ uuid.UUID) ([]*domain.Order, error) {
	out := []*domain.Order{}
	for _, o := range r.orders {
		out = append(out, o)
	}
	return out, nil
}

func (r *memOrderRepo) UpdateStatus(_ context.Context, id uuid.UUID, st domain.OrderStatus) error {
	r.statuses[id] = st
	return nil
}

type memOrderItemRepo struct{}

func (r *memOrderItemRepo) Create(_ context.Context, _ *domain.OrderItem) error { return nil }
func (r *memOrderItemRepo) GetByOrderID(_ context.Context, _ uuid.UUID) ([]*domain.OrderItem, error) {
	return nil, nil
}

var _ repository.OrderRepository = (*memOrderRepo)(nil)
var _ repository.OrderItemRepository = (*memOrderItemRepo)(nil)

// ---- セットアップ ----

type stack struct {
	orderRepo          *memOrderRepo
	usecase            usecase.OrderManagementUsecase
	paymentInspector   *paymenttest.Inspector
	shipInspector      *shippingtest.Inspector
	inventoryInspector *inventorytest.Inspector
	shippingAddr       string
}

func startStack(t *testing.T) *stack {
	t.Helper()

	paymentAddr, stopPayment, paymentInspector, err := paymenttest.StartServer()
	if err != nil {
		t.Fatalf("start payment server: %v", err)
	}
	t.Cleanup(stopPayment)

	shippingAddr, stopShipping, shipInspector, err := shippingtest.StartServer(paymentAddr)
	if err != nil {
		t.Fatalf("start shipping server: %v", err)
	}
	t.Cleanup(stopShipping)

	paymentClient, err := client.NewPaymentClient(paymentAddr)
	if err != nil {
		t.Fatalf("payment client: %v", err)
	}
	t.Cleanup(func() { paymentClient.Close() })

	shippingClient, err := client.NewShippingClient(shippingAddr)
	if err != nil {
		t.Fatalf("shipping client: %v", err)
	}
	t.Cleanup(func() { shippingClient.Close() })

	inventoryAddr, stopInventory, inventoryInspector, err := inventorytest.StartServer()
	if err != nil {
		t.Fatalf("start inventory server: %v", err)
	}
	t.Cleanup(stopInventory)

	inventoryClient, err := client.NewInventoryClient(inventoryAddr)
	if err != nil {
		t.Fatalf("inventory client: %v", err)
	}
	t.Cleanup(func() { inventoryClient.Close() })

	orderRepo := newMemOrderRepo()
	u := usecase.NewOrderManagementUsecaseFull(
		orderRepo, &memOrderItemRepo{}, paymentClient, shippingClient, nil, inventoryClient)

	return &stack{
		orderRepo:          orderRepo,
		usecase:            u,
		paymentInspector:   paymentInspector,
		shipInspector:      shipInspector,
		inventoryInspector: inventoryInspector,
		shippingAddr:       shippingAddr,
	}
}

func testCtx(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// ---- テスト ----

func TestIntegration_CODOrderToDeliveryConfirmsPayment(t *testing.T) {
	ctx := testCtx(t)
	s := startStack(t)

	// 0. 在庫を用意
	productID := uuid.New()
	if err := s.inventoryInspector.SeedStock(productID.String(), 10); err != nil {
		t.Fatalf("seed stock: %v", err)
	}

	// 1. 代引きで注文
	orderID, err := s.usecase.CreateOrder(ctx, usecase.CreateOrderInput{
		CustomerID:    uuid.New(),
		AddressID:     uuid.New(),
		ShippingFee:   500,
		PaymentMethod: usecase.PaymentMethodCashOnDelivery,
		Items: []usecase.OrderItemInput{
			{ProductID: productID, Quantity: 2, UnitPrice: 1200},
		},
	})
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}

	// 在庫が引き当てられ、決済成功で確定している(order → inventory 連携)
	if _, reserved, _ := s.inventoryInspector.Stock(productID.String()); reserved != 2 {
		t.Errorf("reserved stock = %d, want 2", reserved)
	}
	if sts := s.inventoryInspector.ReservationStatuses(orderID.String()); len(sts) != 1 || sts[0] != "CONFIRMED" {
		t.Errorf("reservation statuses = %v, want [CONFIRMED]", sts)
	}

	if st := s.orderRepo.statuses[orderID]; st != domain.OrderStatusConfirmed {
		t.Errorf("order status = %s, want CONFIRMED", st)
	}
	status, method, amount, ok := s.paymentInspector.PaymentByOrder(orderID.String())
	if !ok {
		t.Fatal("payment not created")
	}
	if status != "pending" || method != "cash_on_delivery" {
		t.Errorf("payment = %s/%s, want pending/cash_on_delivery", status, method)
	}
	if amount != 2900 {
		t.Errorf("payment amount = %d, want 2900 (2×1200 + 送料500)", amount)
	}
	shipmentID, ok := s.shipInspector.ShipmentIDByOrder(orderID.String())
	if !ok {
		t.Fatal("shipment not created")
	}

	// 2. 出荷 → 配達完了(配送業者の通知を実 gRPC で模擬)
	conn, err := grpc.NewClient(s.shippingAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial shipping: %v", err)
	}
	defer conn.Close()
	stub := shippingpb.NewShippingServiceClient(conn)

	if _, err := stub.RegisterTrackingNumber(ctx, &shippingpb.RegisterTrackingNumberRequest{
		ShipmentId:     shipmentID,
		TrackingNumber: "TRK-E2E-1",
	}); err != nil {
		t.Fatalf("RegisterTrackingNumber: %v", err)
	}
	resp, err := stub.UpdateShipmentStatus(ctx, &shippingpb.UpdateShipmentStatusRequest{
		ShipmentId: shipmentID,
		NewStatus:  shippingpb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED,
	})
	if err != nil {
		t.Fatalf("UpdateShipmentStatus: %v", err)
	}
	if !resp.GetSuccess() {
		t.Error("expected delivery success")
	}

	// 3. 配達完了によって代引き決済が completed になっている(shipping → payment 連携)
	status, _, _, _ = s.paymentInspector.PaymentByOrder(orderID.String())
	if status != "completed" {
		t.Errorf("payment status after delivery = %s, want completed", status)
	}
}

func TestIntegration_CardOrderCancelRefunds(t *testing.T) {
	ctx := testCtx(t)
	s := startStack(t)

	productID := uuid.New()
	if err := s.inventoryInspector.SeedStock(productID.String(), 3); err != nil {
		t.Fatalf("seed stock: %v", err)
	}

	// 1. クレジットカードで注文(即時決済)
	orderID, err := s.usecase.CreateOrder(ctx, usecase.CreateOrderInput{
		CustomerID:      uuid.New(),
		AddressID:       uuid.New(),
		ShippingFee:     500,
		PaymentMethodID: "pm_test",
		Items: []usecase.OrderItemInput{
			{ProductID: productID, Quantity: 1, UnitPrice: 5000},
		},
	})
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}
	if st := s.orderRepo.statuses[orderID]; st != domain.OrderStatusPaid {
		t.Errorf("order status = %s, want PAID", st)
	}
	status, _, _, ok := s.paymentInspector.PaymentByOrder(orderID.String())
	if !ok || status != "completed" {
		t.Fatalf("payment should be completed, got %s (ok=%v)", status, ok)
	}

	// 2. キャンセル → 自動返金
	if err := s.usecase.CancelOrder(ctx, orderID); err != nil {
		t.Fatalf("CancelOrder: %v", err)
	}
	if st := s.orderRepo.statuses[orderID]; st != domain.OrderStatusCancelled {
		t.Errorf("order status = %s, want CANCELLED", st)
	}
	status, _, _, _ = s.paymentInspector.PaymentByOrder(orderID.String())
	if status != "refunded" {
		t.Errorf("payment status = %s, want refunded", status)
	}
	if n := s.paymentInspector.RefundCount(); n != 1 {
		t.Fatalf("expected 1 refund record, got %d", n)
	}
	if amounts := s.paymentInspector.RefundAmounts(); len(amounts) != 1 || amounts[0] != 5500 {
		t.Errorf("refund amounts = %v, want [5500] (全額)", amounts)
	}
}

func TestIntegration_CancelWithoutPaymentIsNoop(t *testing.T) {
	ctx := testCtx(t)
	s := startStack(t)

	// 決済のない注文のキャンセル(payment が NotFound を返し、返金はスキップされる)
	orderID := uuid.New()
	if err := s.usecase.CancelOrder(ctx, orderID); err != nil {
		t.Fatalf("CancelOrder should succeed without payment: %v", err)
	}
	if n := s.paymentInspector.RefundCount(); n != 0 {
		t.Errorf("no refunds expected, got %d", n)
	}
}

func TestIntegration_InsufficientStockFailsOrderBeforePayment(t *testing.T) {
	ctx := testCtx(t)
	s := startStack(t)

	productID := uuid.New()
	if err := s.inventoryInspector.SeedStock(productID.String(), 1); err != nil {
		t.Fatalf("seed stock: %v", err)
	}

	// 在庫1に対して5個注文 → 引当が FailedPrecondition で失敗し、決済まで進まない
	_, err := s.usecase.CreateOrder(ctx, usecase.CreateOrderInput{
		CustomerID:  uuid.New(),
		AddressID:   uuid.New(),
		ShippingFee: 500,
		Items: []usecase.OrderItemInput{
			{ProductID: productID, Quantity: 5, UnitPrice: 1000},
		},
	})
	if err == nil {
		t.Fatal("expected error when stock is insufficient")
	}
	for id := range s.orderRepo.orders {
		if status, _, _, ok := s.paymentInspector.PaymentByOrder(id.String()); ok {
			t.Errorf("payment should not exist for failed order, got %s", status)
		}
		if st := s.orderRepo.statuses[id]; st != domain.OrderStatusCancelled {
			t.Errorf("order status = %s, want CANCELLED", st)
		}
	}
	// 引当は残っていない(補償で解放済み or そもそも引当されていない)
	if _, reserved, _ := s.inventoryInspector.Stock(productID.String()); reserved != 0 {
		t.Errorf("reserved stock = %d, want 0", reserved)
	}
}
