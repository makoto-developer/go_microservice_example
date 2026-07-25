// Package testsupport は他サービスの統合テスト向けに、
// インメモリ(DB なし)の inventory gRPC サーバを起動するヘルパーを提供する。
// 本番コードから import してはならない。
package testsupport

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/domain"
	handlergrpc "github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/handler/grpc"
	"github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/inventory/proto"
)

// Inspector はテストからサーバ内部の状態を覗き、在庫を仕込むための口。
type Inspector struct {
	inventories  *memInventoryRepo
	reservations *memReservationRepo
}

// SeedStock は商品の在庫を登録する(テストの前提データ)。
func (i *Inspector) SeedStock(productID string, quantity int) error {
	pid, err := uuid.Parse(productID)
	if err != nil {
		return err
	}
	now := time.Now()
	inv := &domain.Inventory{
		ID:        uuid.New(),
		ProductID: pid,
		ShopID:    uuid.New(),
		Quantity:  quantity,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return i.inventories.Create(context.Background(), inv)
}

// Stock は商品の (総数, 引当済み) を返す。
func (i *Inspector) Stock(productID string) (quantity, reserved int, ok bool) {
	pid, err := uuid.Parse(productID)
	if err != nil {
		return 0, 0, false
	}
	inv, _ := i.inventories.GetByProductID(context.Background(), pid, nil)
	if inv == nil {
		return 0, 0, false
	}
	return inv.Quantity, inv.ReservedQuantity, true
}

// ReservationStatuses は注文に紐づく引当の状態一覧を返す(状態を問わず全件)。
func (i *Inspector) ReservationStatuses(orderID string) []string {
	oid, err := uuid.Parse(orderID)
	if err != nil {
		return nil
	}
	out := []string{}
	for _, r := range i.reservations.items {
		if r.OrderID == oid {
			out = append(out, string(r.Status))
		}
	}
	return out
}

// StartServer はインメモリの inventory gRPC サーバを起動する。
func StartServer() (addr string, stop func(), inspect *Inspector, err error) {
	inventories := newMemInventoryRepo()
	reservations := newMemReservationRepo()
	handler := handlergrpc.NewInventoryServiceHandler(
		usecase.NewInventoryManagementUsecase(inventories, reservations),
	)
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, nil, err
	}
	server := grpc.NewServer()
	pb.RegisterInventoryServiceServer(server, handler)
	go server.Serve(lis)
	return lis.Addr().String(), server.Stop, &Inspector{inventories: inventories, reservations: reservations}, nil
}

// ---- インメモリリポジトリ ----

type memInventoryRepo struct {
	items map[uuid.UUID]*domain.Inventory
}

func newMemInventoryRepo() *memInventoryRepo {
	return &memInventoryRepo{items: map[uuid.UUID]*domain.Inventory{}}
}

func (r *memInventoryRepo) Create(_ context.Context, inv *domain.Inventory) error {
	r.items[inv.ID] = inv
	return nil
}

func (r *memInventoryRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Inventory, error) {
	if inv, ok := r.items[id]; ok {
		return inv, nil
	}
	return nil, fmt.Errorf("inventory %s not found", id)
}

func (r *memInventoryRepo) GetByProductID(_ context.Context, productID uuid.UUID, _ *uuid.UUID) (*domain.Inventory, error) {
	for _, inv := range r.items {
		if inv.ProductID == productID {
			return inv, nil
		}
	}
	return nil, fmt.Errorf("inventory not found for product %s", productID)
}

func (r *memInventoryRepo) Update(_ context.Context, inv *domain.Inventory) error {
	r.items[inv.ID] = inv
	return nil
}

func (r *memInventoryRepo) UpdateQuantity(_ context.Context, id uuid.UUID, quantity int) error {
	if inv, ok := r.items[id]; ok {
		inv.Quantity = quantity
	}
	return nil
}

func (r *memInventoryRepo) Reserve(_ context.Context, id uuid.UUID, quantity int) error {
	inv, ok := r.items[id]
	if !ok {
		return fmt.Errorf("inventory %s not found", id)
	}
	if !inv.CanReserve(quantity) {
		return fmt.Errorf("insufficient stock: available %d, requested %d", inv.AvailableQuantity(), quantity)
	}
	inv.ReservedQuantity += quantity
	return nil
}

func (r *memInventoryRepo) Release(_ context.Context, id uuid.UUID, quantity int) error {
	if inv, ok := r.items[id]; ok {
		inv.ReservedQuantity -= quantity
		if inv.ReservedQuantity < 0 {
			inv.ReservedQuantity = 0
		}
	}
	return nil
}

type memReservationRepo struct {
	items map[uuid.UUID]*domain.Reservation
}

func newMemReservationRepo() *memReservationRepo {
	return &memReservationRepo{items: map[uuid.UUID]*domain.Reservation{}}
}

func (r *memReservationRepo) Create(_ context.Context, res *domain.Reservation) error {
	r.items[res.ID] = res
	return nil
}

func (r *memReservationRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Reservation, error) {
	return r.items[id], nil
}

func (r *memReservationRepo) GetByOrderID(_ context.Context, orderID uuid.UUID) ([]*domain.Reservation, error) {
	out := []*domain.Reservation{}
	for _, res := range r.items {
		if res.OrderID == orderID && res.Status == domain.ReservationStatusPending {
			out = append(out, res)
		}
	}
	return out, nil
}

func (r *memReservationRepo) UpdateStatus(_ context.Context, id uuid.UUID, status domain.ReservationStatus) error {
	if res, ok := r.items[id]; ok {
		res.Status = status
	}
	return nil
}

func (r *memReservationRepo) GetExpiredPending(_ context.Context) ([]*domain.Reservation, error) {
	out := []*domain.Reservation{}
	for _, res := range r.items {
		if res.Status == domain.ReservationStatusPending && res.ExpiresAt.Before(time.Now()) {
			out = append(out, res)
		}
	}
	return out, nil
}

func (r *memReservationRepo) DeleteExpired(_ context.Context) error { return nil }
