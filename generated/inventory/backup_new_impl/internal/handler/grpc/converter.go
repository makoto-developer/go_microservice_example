package grpc

import (
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/inventory/internal/domain"
	pb "github.com/makoto-developer/go_microservice_example/proto/inventory-service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func domainInventoryToProto(inv *domain.Inventory) *pb.Inventory {
	return &pb.Inventory{
		Id:               inv.ID.String(),
		ProductId:        inv.ProductID.String(),
		ShopId:           inv.ShopID.String(),
		Quantity:         int32(inv.Quantity),
		ReservedQuantity: int32(inv.ReservedQuantity),
		CreatedAt:        timestampProto(inv.CreatedAt),
		UpdatedAt:        timestampProto(inv.UpdatedAt),
	}
}

func timestampProto(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
