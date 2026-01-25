package grpc

import (
	pb "github.com/makoto-developer/go_microservice_example/generated/shop/proto/shop_service/v1"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
	"fmt"
)

// convertToProtoShopStatus converts domain ShopStatus to proto ShopStatus
func convertToProtoShopStatus(status domain.ShopStatus) pb.ShopStatus {
	switch status {
	case domain.ShopStatusPendingApproval:
		return pb.ShopStatus_PENDING_APPROVAL
	case domain.ShopStatusApproved:
		return pb.ShopStatus_APPROVED
	case domain.ShopStatusSuspended:
		return pb.ShopStatus_SUSPENDED
	default:
		return pb.ShopStatus_SHOP_STATUS_UNSPECIFIED
	}
}

// convertToProtoShop converts domain Shop to proto Shop
func convertToProtoShop(s *domain.Shop) *pb.Shop {
	shop := &pb.Shop{
		Id:            s.ID.String(),
		OwnerId:       s.OwnerID.String(),
		Name:          s.Name,
		Description:   s.Description,
		LogoUrl:       s.LogoImageURL,
		OwnerName:     s.OwnerName,
		PhoneNumber:   s.OwnerPhone,
		BusinessHours: s.BusinessHours,
		ReturnPolicy:  s.ReturnPolicy,
		Status:        convertToProtoShopStatus(s.Status),
		Published:     s.IsPublic,
		CreatedAt:     timestamppb.New(s.CreatedAt),
		UpdatedAt:     timestamppb.New(s.UpdatedAt),
	}

	return shop
}

// convertToProtoProduct converts domain Product to proto Product
func convertToProtoProduct(p *domain.Product) *pb.Product {
	product := &pb.Product{
		Id:          p.ID.String(),
		ShopId:      p.ShopID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       fmt.Sprintf("%d", p.Price),
		Category:    p.CategoryID.String(),
		StockQuantity: int32(p.StockCount),
		Published:   p.IsPublic,
		Deleted:     p.IsDeleted,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}

	if p.Weight != nil {
		product.Weight = fmt.Sprintf("%.2f", *p.Weight)
	}

	if p.Dimensions != nil {
		product.Size = *p.Dimensions
	}

	if p.JANCode != nil {
		product.JanCode = *p.JANCode
	}

	return product
}
