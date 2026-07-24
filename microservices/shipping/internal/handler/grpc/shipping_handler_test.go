package grpc

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/shipping/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// --- フェイク ---

type fakeShipmentRepo struct {
	shipments map[uuid.UUID]*domain.Shipment
}

func newFakeShipmentRepo() *fakeShipmentRepo {
	return &fakeShipmentRepo{shipments: map[uuid.UUID]*domain.Shipment{}}
}

func (r *fakeShipmentRepo) Create(_ context.Context, s *domain.Shipment) error {
	r.shipments[s.ID] = s
	return nil
}

func (r *fakeShipmentRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Shipment, error) {
	return r.shipments[id], nil
}

func (r *fakeShipmentRepo) GetByOrderID(_ context.Context, orderID uuid.UUID) (*domain.Shipment, error) {
	for _, s := range r.shipments {
		if s.OrderID == orderID {
			return s, nil
		}
	}
	return nil, nil
}

func (r *fakeShipmentRepo) GetByCustomerID(_ context.Context, customerID uuid.UUID) ([]*domain.Shipment, error) {
	out := []*domain.Shipment{}
	for _, s := range r.shipments {
		if s.CustomerID == customerID {
			out = append(out, s)
		}
	}
	return out, nil
}

func (r *fakeShipmentRepo) UpdateStatus(_ context.Context, id uuid.UUID, st domain.ShipmentStatus) error {
	if s, ok := r.shipments[id]; ok {
		s.Status = st
	}
	return nil
}

func (r *fakeShipmentRepo) UpdateTracking(_ context.Context, id uuid.UUID, tn string) error {
	if s, ok := r.shipments[id]; ok {
		s.TrackingNumber = tn
	}
	return nil
}

type fakePaymentClient struct {
	confirmedOrders []string
	err             error
}

func (c *fakePaymentClient) ConfirmCODByOrder(_ context.Context, orderID string) error {
	c.confirmedOrders = append(c.confirmedOrders, orderID)
	return c.err
}

func (c *fakePaymentClient) Close() error { return nil }

func newTestHandler() (*ShippingServiceHandler, *fakeShipmentRepo, *fakePaymentClient) {
	repo := newFakeShipmentRepo()
	payment := &fakePaymentClient{}
	h := NewShippingServiceHandler(usecase.NewCreateShipmentUsecase(repo), repo, payment)
	return h, repo, payment
}

// --- テスト ---

func TestShipmentLifecycle_DeliveredConfirmsCOD(t *testing.T) {
	h, repo, payment := newTestHandler()
	orderID := uuid.New()

	created, err := h.CreateShipment(context.Background(), &pb.CreateShipmentRequest{
		OrderId:         orderID.String(),
		CustomerId:      uuid.New().String(),
		ShippingAddress: "東京都渋谷区1-2-3",
	})
	if err != nil {
		t.Fatalf("CreateShipment: %v", err)
	}
	shipmentID := created.GetShipmentId()

	// 追跡番号の登録 → shipped
	if _, err := h.RegisterTrackingNumber(context.Background(), &pb.RegisterTrackingNumberRequest{
		ShipmentId:     shipmentID,
		TrackingNumber: "TRK-0001",
	}); err != nil {
		t.Fatalf("RegisterTrackingNumber: %v", err)
	}
	sid := uuid.MustParse(shipmentID)
	if repo.shipments[sid].Status != domain.ShipmentStatusShipped {
		t.Errorf("status = %s, want shipped", repo.shipments[sid].Status)
	}

	// 配達完了 → COD 集金確定が呼ばれる
	resp, err := h.UpdateShipmentStatus(context.Background(), &pb.UpdateShipmentStatusRequest{
		ShipmentId: shipmentID,
		NewStatus:  pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED,
	})
	if err != nil {
		t.Fatalf("UpdateShipmentStatus: %v", err)
	}
	if !resp.GetSuccess() {
		t.Error("expected success")
	}
	if len(payment.confirmedOrders) != 1 || payment.confirmedOrders[0] != orderID.String() {
		t.Errorf("COD confirmation calls = %v, want [%s]", payment.confirmedOrders, orderID)
	}

	// 配達完了後は更新できない
	_, err = h.UpdateShipmentStatus(context.Background(), &pb.UpdateShipmentStatusRequest{
		ShipmentId: shipmentID,
		NewStatus:  pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT,
	})
	if status.Code(err) != codes.FailedPrecondition {
		t.Errorf("expected FailedPrecondition after delivery, got %v", err)
	}
}

func TestUpdateShipmentStatus_NonDeliveredDoesNotTouchPayment(t *testing.T) {
	h, _, payment := newTestHandler()

	created, err := h.CreateShipment(context.Background(), &pb.CreateShipmentRequest{
		OrderId:         uuid.New().String(),
		ShippingAddress: "大阪府大阪市4-5-6",
	})
	if err != nil {
		t.Fatalf("CreateShipment: %v", err)
	}

	if _, err := h.UpdateShipmentStatus(context.Background(), &pb.UpdateShipmentStatusRequest{
		ShipmentId: created.GetShipmentId(),
		NewStatus:  pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT,
	}); err != nil {
		t.Fatalf("UpdateShipmentStatus: %v", err)
	}
	if len(payment.confirmedOrders) != 0 {
		t.Errorf("payment should not be called before delivery, got %v", payment.confirmedOrders)
	}
}

func TestCalculateShippingFee_Table(t *testing.T) {
	h, _, _ := newTestHandler()

	cases := []struct {
		method string
		weight int32
		want   string
	}{
		{"standard", 500, "500"},
		{"express", 500, "1000"},
		{"standard", 15000, "1000"},
	}
	for _, c := range cases {
		resp, err := h.CalculateShippingFee(context.Background(), &pb.CalculateShippingFeeRequest{
			ShippingMethod:   c.method,
			TotalWeightGrams: c.weight,
		})
		if err != nil {
			t.Fatalf("CalculateShippingFee(%s): %v", c.method, err)
		}
		if resp.GetFee() != c.want {
			t.Errorf("fee(%s, %dg) = %s, want %s", c.method, c.weight, resp.GetFee(), c.want)
		}
	}
}

func TestGetShipmentByOrder_NotFound(t *testing.T) {
	h, _, _ := newTestHandler()

	_, err := h.GetShipmentByOrder(context.Background(), &pb.GetShipmentByOrderRequest{
		OrderId: uuid.New().String(),
	})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound, got %v", err)
	}
}
