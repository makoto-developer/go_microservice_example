// Package testsupport は他サービスの統合テスト向けに、
// インメモリ(DB なし)の payment gRPC サーバを起動するヘルパーを提供する。
// 本番コードから import してはならない。
package testsupport

import (
	"context"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/domain"
	handlergrpc "github.com/makoto-developer/go_microservice_example/microservices/payment/internal/handler/grpc"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/repository"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/payment/proto"
)

// Inspector はテストからサーバ内部の状態を覗くための読み取り口。
type Inspector struct {
	payments *memPaymentRepo
	refunds  *memRefundRepo
}

// PaymentByOrder は注文に紐づく決済の (status, method, amount) を返す。
func (i *Inspector) PaymentByOrder(orderID string) (status string, method string, amount int, ok bool) {
	id, err := uuid.Parse(orderID)
	if err != nil {
		return "", "", 0, false
	}
	p, _ := i.payments.GetByOrderID(context.Background(), id)
	if p == nil {
		return "", "", 0, false
	}
	return string(p.Status), string(p.PaymentMethod), p.Amount, true
}

// RefundCount は作成された返金レコード数を返す。
func (i *Inspector) RefundCount() int { return len(i.refunds.refunds) }

// RefundAmounts は返金金額の一覧を返す。
func (i *Inspector) RefundAmounts() []int {
	out := []int{}
	for _, r := range i.refunds.refunds {
		out = append(out, r.Amount)
	}
	return out
}

// StartServer はインメモリの payment gRPC サーバを 127.0.0.1 の空きポートで起動する。
func StartServer() (addr string, stop func(), inspect *Inspector, err error) {
	payments := newMemPaymentRepo()
	refunds := newMemRefundRepo()
	handler := handlergrpc.NewPaymentServiceHandler(
		usecase.NewProcessPaymentUsecase(payments),
		usecase.NewCODPaymentUsecase(payments),
		usecase.NewRefundPaymentUsecase(payments, refunds),
		payments,
		refunds,
	)
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, nil, err
	}
	server := grpc.NewServer()
	pb.RegisterPaymentServiceServer(server, handler)
	go server.Serve(lis)
	return lis.Addr().String(), server.Stop, &Inspector{payments: payments, refunds: refunds}, nil
}

// ---- インメモリリポジトリ ----

type memPaymentRepo struct{ payments map[uuid.UUID]*domain.Payment }

func newMemPaymentRepo() *memPaymentRepo {
	return &memPaymentRepo{payments: map[uuid.UUID]*domain.Payment{}}
}

func (r *memPaymentRepo) Create(_ context.Context, p *domain.Payment) error {
	cp := *p
	r.payments[p.ID] = &cp
	return nil
}

func (r *memPaymentRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Payment, error) {
	if p, ok := r.payments[id]; ok {
		cp := *p
		return &cp, nil
	}
	return nil, nil
}

func (r *memPaymentRepo) GetByOrderID(_ context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	for _, p := range r.payments {
		if p.OrderID == orderID {
			cp := *p
			return &cp, nil
		}
	}
	return nil, nil
}

func (r *memPaymentRepo) List(_ context.Context, filter repository.PaymentListFilter) ([]*domain.Payment, int, error) {
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

func (r *memPaymentRepo) UpdateStatus(_ context.Context, id uuid.UUID, s domain.PaymentStatus, txid string) error {
	if p, ok := r.payments[id]; ok {
		p.Status = s
		p.TransactionID = txid
	}
	return nil
}

type memRefundRepo struct{ refunds map[uuid.UUID]*domain.Refund }

func newMemRefundRepo() *memRefundRepo {
	return &memRefundRepo{refunds: map[uuid.UUID]*domain.Refund{}}
}

func (r *memRefundRepo) Create(_ context.Context, refund *domain.Refund) error {
	cp := *refund
	r.refunds[refund.ID] = &cp
	return nil
}

func (r *memRefundRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Refund, error) {
	if refund, ok := r.refunds[id]; ok {
		cp := *refund
		return &cp, nil
	}
	return nil, nil
}
