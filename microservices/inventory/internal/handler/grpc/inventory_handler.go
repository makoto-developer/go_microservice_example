package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/inventory/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InventoryServiceHandler struct {
	pb.UnimplementedInventoryServiceServer
	inventoryMgmt usecase.InventoryManagementUsecase
}

func NewInventoryServiceHandler(inventoryMgmt usecase.InventoryManagementUsecase) *InventoryServiceHandler {
	return &InventoryServiceHandler{inventoryMgmt: inventoryMgmt}
}

func (h *InventoryServiceHandler) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
	inventoryID, err := uuid.Parse(req.InventoryId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid inventory_id: %v", err)
	}

	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}

	// For now, use the inventory_id directly
	// In a real implementation, you'd fetch the inventory to get product/variation IDs
	input := usecase.ReserveInventoryInput{
		ProductID:   inventoryID, // Simplified for now
		VariationID: nil,
		OrderID:     orderID,
		Quantity:    int(req.Quantity),
	}

	_, err = h.inventoryMgmt.ReserveInventory(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reserve: %v", err)
	}

	return &pb.ReserveStockResponse{}, nil
}

func (h *InventoryServiceHandler) ReleaseStock(ctx context.Context, req *pb.ReleaseStockRequest) (*pb.ReleaseStockResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}

	if err := h.inventoryMgmt.ReleaseInventory(ctx, orderID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to release: %v", err)
	}

	return &pb.ReleaseStockResponse{}, nil
}

func inventoryToProto(inv *domain.Inventory) *pb.Inventory {
	out := &pb.Inventory{
		Id:                inv.ID.String(),
		ProductId:         inv.ProductID.String(),
		ShopId:            inv.ShopID.String(),
		Quantity:          int32(inv.Quantity),
		ReservedQuantity:  int32(inv.ReservedQuantity),
		AvailableQuantity: int32(inv.AvailableQuantity()),
		CreatedAt:         timestamppb.New(inv.CreatedAt),
		UpdatedAt:         timestamppb.New(inv.UpdatedAt),
	}
	if inv.VariationID != nil {
		out.VariationId = inv.VariationID.String()
	}
	return out
}

// RegisterInventory は商品の在庫レコードを新規登録する(商品登録時に呼ばれる)。
func (h *InventoryServiceHandler) RegisterInventory(ctx context.Context, req *pb.RegisterInventoryRequest) (*pb.RegisterInventoryResponse, error) {
	productID, err := uuid.Parse(req.GetProductId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product_id: %v", err)
	}
	shopID, err := uuid.Parse(req.GetShopId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid shop_id: %v", err)
	}
	var variationID *uuid.UUID
	if req.GetVariationId() != "" {
		vid, err := uuid.Parse(req.GetVariationId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid variation_id: %v", err)
		}
		variationID = &vid
	}

	inv, err := h.inventoryMgmt.RegisterInventory(ctx, productID, variationID, shopID, int(req.GetInitialQuantity()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register inventory: %v", err)
	}
	return &pb.RegisterInventoryResponse{Inventory: inventoryToProto(inv), Message: "Inventory registered"}, nil
}

// UpdateInventoryQuantity は入荷・調整などの在庫数変更を反映する。
func (h *InventoryServiceHandler) UpdateInventoryQuantity(ctx context.Context, req *pb.UpdateInventoryQuantityRequest) (*pb.UpdateInventoryQuantityResponse, error) {
	inventoryID, err := uuid.Parse(req.GetInventoryId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid inventory_id: %v", err)
	}
	current, err := h.inventoryMgmt.GetInventory(ctx, inventoryID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "inventory not found: %v", err)
	}

	newQuantity := current.Quantity + int(req.GetChangeQuantity())
	inv, err := h.inventoryMgmt.AdjustQuantity(ctx, inventoryID, newQuantity)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to adjust quantity: %v", err)
	}
	return &pb.UpdateInventoryQuantityResponse{
		Inventory: inventoryToProto(inv),
		Message:   "Quantity updated",
	}, nil
}

// GetInventory は在庫を ID で取得する。
func (h *InventoryServiceHandler) GetInventory(ctx context.Context, req *pb.GetInventoryRequest) (*pb.GetInventoryResponse, error) {
	inventoryID, err := uuid.Parse(req.GetInventoryId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid inventory_id: %v", err)
	}
	inv, err := h.inventoryMgmt.GetInventory(ctx, inventoryID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "inventory not found: %v", err)
	}
	return &pb.GetInventoryResponse{Inventory: inventoryToProto(inv)}, nil
}

// GetInventoryByProduct は商品 ID から在庫を引く(商品ページの在庫表示用)。
func (h *InventoryServiceHandler) GetInventoryByProduct(ctx context.Context, req *pb.GetInventoryByProductRequest) (*pb.GetInventoryByProductResponse, error) {
	productID, err := uuid.Parse(req.GetProductId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product_id: %v", err)
	}
	var variationID *uuid.UUID
	if req.GetVariationId() != "" {
		vid, err := uuid.Parse(req.GetVariationId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid variation_id: %v", err)
		}
		variationID = &vid
	}
	inv, err := h.inventoryMgmt.GetInventoryByProduct(ctx, productID, variationID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "inventory not found: %v", err)
	}
	return &pb.GetInventoryByProductResponse{Inventory: inventoryToProto(inv)}, nil
}

// BulkGetInventory は複数商品の在庫をまとめて返す(一覧画面の在庫列用)。
// 在庫レコードが無い商品はスキップする。
func (h *InventoryServiceHandler) BulkGetInventory(ctx context.Context, req *pb.BulkGetInventoryRequest) (*pb.BulkGetInventoryResponse, error) {
	out := []*pb.Inventory{}
	for _, pid := range req.GetProductIds() {
		productID, err := uuid.Parse(pid)
		if err != nil {
			continue
		}
		inv, err := h.inventoryMgmt.GetInventoryByProduct(ctx, productID, nil)
		if err != nil || inv == nil {
			continue
		}
		out = append(out, inventoryToProto(inv))
	}
	return &pb.BulkGetInventoryResponse{Inventories: out}, nil
}

// BulkReserveStock は注文の全商品の在庫を一括で引き当てる(order サービスから呼ばれる)。
// いずれかの引当に失敗したら、その注文で引当済みの分を解放してエラーを返す(全部成功 or 全部なし)。
func (h *InventoryServiceHandler) BulkReserveStock(ctx context.Context, req *pb.BulkReserveStockRequest) (*pb.BulkReserveStockResponse, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}
	if len(req.GetReservations()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "reservations are required")
	}

	for _, r := range req.GetReservations() {
		// 注: このサンプルでは inventory_id フィールドに product_id を渡す運用(ReserveStock と同じ簡略化)
		productID, err := uuid.Parse(r.GetInventoryId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid inventory_id: %v", err)
		}
		input := usecase.ReserveInventoryInput{
			ProductID: productID,
			OrderID:   orderID,
			Quantity:  int(r.GetQuantity()),
		}
		if _, err := h.inventoryMgmt.ReserveInventory(ctx, input); err != nil {
			// 途中まで引き当てた分を解放(補償)
			if relErr := h.inventoryMgmt.ReleaseInventory(ctx, orderID); relErr != nil {
				return nil, status.Errorf(codes.Internal, "reserve failed (%v) and rollback also failed: %v", err, relErr)
			}
			return nil, status.Errorf(codes.FailedPrecondition, "failed to reserve stock: %v", err)
		}
	}

	return &pb.BulkReserveStockResponse{Message: "Stock reserved"}, nil
}

func (h *InventoryServiceHandler) ConfirmStock(ctx context.Context, req *pb.ConfirmStockRequest) (*pb.ConfirmStockResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}

	if err := h.inventoryMgmt.ConfirmInventory(ctx, orderID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to confirm: %v", err)
	}

	return &pb.ConfirmStockResponse{}, nil
}

// ReleaseExpiredReservations は期限切れの引当を解放する(在庫バッチから定期実行される)。
func (h *InventoryServiceHandler) ReleaseExpiredReservations(ctx context.Context, req *pb.ReleaseExpiredReservationsRequest) (*pb.ReleaseExpiredReservationsResponse, error) {
	if err := h.inventoryMgmt.ReleaseExpiredReservations(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to release expired reservations: %v", err)
	}
	return &pb.ReleaseExpiredReservationsResponse{Message: "Expired reservations released"}, nil
}

// CheckStockAlert は在庫が閾値を割っていないかを返す(在庫バッチが巡回する)。
// このサンプルでは閾値はテーブルに持たず固定値 10 とする。
func (h *InventoryServiceHandler) CheckStockAlert(ctx context.Context, req *pb.CheckStockAlertRequest) (*pb.CheckStockAlertResponse, error) {
	const alertThreshold = 10
	inventoryID, err := uuid.Parse(req.GetInventoryId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid inventory_id: %v", err)
	}
	inv, err := h.inventoryMgmt.GetInventory(ctx, inventoryID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "inventory not found: %v", err)
	}
	return &pb.CheckStockAlertResponse{
		IsLowStock:      inv.AvailableQuantity() < alertThreshold,
		CurrentQuantity: int32(inv.AvailableQuantity()),
		AlertThreshold:  alertThreshold,
	}, nil
}

// GetInventoryHistory は在庫変動履歴を返す。
// このサンプルは履歴テーブルを持たないため常に空を返す(正常応答)。
func (h *InventoryServiceHandler) GetInventoryHistory(ctx context.Context, req *pb.GetInventoryHistoryRequest) (*pb.GetInventoryHistoryResponse, error) {
	return &pb.GetInventoryHistoryResponse{
		History:  []*pb.InventoryHistory{},
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
	}, nil
}

// RecordStockTaking は棚卸しの実数を記録し、システム在庫を実数に合わせる。
func (h *InventoryServiceHandler) RecordStockTaking(ctx context.Context, req *pb.RecordStockTakingRequest) (*pb.RecordStockTakingResponse, error) {
	inventoryID, err := uuid.Parse(req.GetInventoryId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid inventory_id: %v", err)
	}
	before, err := h.inventoryMgmt.GetInventory(ctx, inventoryID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "inventory not found: %v", err)
	}

	actual := int(req.GetActualQuantity())
	if _, err := h.inventoryMgmt.AdjustQuantity(ctx, inventoryID, actual); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to record stock taking: %v", err)
	}

	return &pb.RecordStockTakingResponse{
		StockTaking: &pb.StockTaking{
			Id:               uuid.New().String(),
			InventoryId:      req.GetInventoryId(),
			ShopId:           req.GetShopId(),
			SystemQuantity:   int32(before.Quantity),
			ActualQuantity:   int32(actual),
			Difference:       int32(actual - before.Quantity),
			DifferenceReason: req.GetDifferenceReason(),
			Operator:         req.GetOperator(),
			CreatedAt:        timestamppb.Now(),
		},
		Message: "Stock taking recorded",
	}, nil
}

// GetStockTakingHistory は棚卸し履歴を返す。
// このサンプルは履歴テーブルを持たないため常に空を返す(正常応答)。
func (h *InventoryServiceHandler) GetStockTakingHistory(ctx context.Context, req *pb.GetStockTakingHistoryRequest) (*pb.GetStockTakingHistoryResponse, error) {
	return &pb.GetStockTakingHistoryResponse{History: []*pb.StockTaking{}}, nil
}
