// Package testsupport は他サービスの統合テスト向けに、
// インメモリ(DB なし)の shipping gRPC サーバを起動するヘルパーを提供する。
// 本番コードから import してはならない。
package testsupport

import (
	"context"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/client"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/domain"
	handlergrpc "github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/handler/grpc"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/shipping/proto"
)

// Inspector はテストからサーバ内部の状態を覗くための読み取り口。
type Inspector struct {
	shipments *memShipmentRepo
}

// ShipmentIDByOrder は注文に紐づく出荷 ID を返す。
func (i *Inspector) ShipmentIDByOrder(orderID string) (string, bool) {
	id, err := uuid.Parse(orderID)
	if err != nil {
		return "", false
	}
	s, _ := i.shipments.GetByOrderID(context.Background(), id)
	if s == nil {
		return "", false
	}
	return s.ID.String(), true
}

// StartServer はインメモリの shipping gRPC サーバを起動する。
// paymentAddr が空でなければ、配達完了時に代引きの集金確定を通知する。
func StartServer(paymentAddr string) (addr string, stop func(), inspect *Inspector, err error) {
	shipments := newMemShipmentRepo()

	var paymentClient client.PaymentClient
	if paymentAddr != "" {
		paymentClient, err = client.NewPaymentClient(paymentAddr)
		if err != nil {
			return "", nil, nil, err
		}
	}

	handler := handlergrpc.NewShippingServiceHandler(
		usecase.NewCreateShipmentUsecase(shipments),
		shipments,
		paymentClient,
	)
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, nil, err
	}
	server := grpc.NewServer()
	pb.RegisterShippingServiceServer(server, handler)
	go server.Serve(lis)
	stop = func() {
		server.Stop()
		if paymentClient != nil {
			paymentClient.Close()
		}
	}
	return lis.Addr().String(), stop, &Inspector{shipments: shipments}, nil
}

// ---- インメモリリポジトリ ----

type memShipmentRepo struct {
	shipments map[uuid.UUID]*domain.Shipment
}

func newMemShipmentRepo() *memShipmentRepo {
	return &memShipmentRepo{shipments: map[uuid.UUID]*domain.Shipment{}}
}

func (r *memShipmentRepo) Create(_ context.Context, s *domain.Shipment) error {
	r.shipments[s.ID] = s
	return nil
}

func (r *memShipmentRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Shipment, error) {
	return r.shipments[id], nil
}

func (r *memShipmentRepo) GetByOrderID(_ context.Context, orderID uuid.UUID) (*domain.Shipment, error) {
	for _, s := range r.shipments {
		if s.OrderID == orderID {
			return s, nil
		}
	}
	return nil, nil
}

func (r *memShipmentRepo) GetByCustomerID(_ context.Context, customerID uuid.UUID) ([]*domain.Shipment, error) {
	out := []*domain.Shipment{}
	for _, s := range r.shipments {
		if s.CustomerID == customerID {
			out = append(out, s)
		}
	}
	return out, nil
}

func (r *memShipmentRepo) UpdateStatus(_ context.Context, id uuid.UUID, st domain.ShipmentStatus) error {
	if s, ok := r.shipments[id]; ok {
		s.Status = st
	}
	return nil
}

func (r *memShipmentRepo) UpdateTracking(_ context.Context, id uuid.UUID, tn string) error {
	if s, ok := r.shipments[id]; ok {
		s.TrackingNumber = tn
	}
	return nil
}
