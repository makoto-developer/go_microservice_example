package grpc

import (
	"strconv"

	pb "github.com/makoto-developer/go_microservice_example/proto/shop_service/v1"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
)

func domainStatusToProto(status domain.ShopStatus) pb.ShopStatus {
	switch status {
	case domain.ShopStatusPendingApproval:
		return pb.ShopStatus_PENDING_APPROVAL
	case domain.ShopStatusApproved:
		return pb.ShopStatus_APPROVED
	case domain.ShopStatusRejected:
		return pb.ShopStatus_REJECTED
	case domain.ShopStatusSuspended:
		return pb.ShopStatus_SUSPENDED
	default:
		return pb.ShopStatus_SHOP_STATUS_UNSPECIFIED
	}
}

func parsePrice(priceStr string) (float64, error) {
	if priceStr == "" {
		return 0, nil
	}
	return strconv.ParseFloat(priceStr, 64)
}

func parseWeight(weightStr string) (*float64, error) {
	if weightStr == "" {
		return nil, nil
	}
	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		return nil, err
	}
	return &weight, nil
}

func formatPrice(price float64) string {
	return strconv.FormatFloat(price, 'f', 2, 64)
}

func formatWeight(weight *float64) string {
	if weight == nil {
		return ""
	}
	return strconv.FormatFloat(*weight, 'f', 2, 64)
}

func domainOrderStatusToProto(status domain.OrderStatus) pb.OrderStatus {
	switch status {
	case domain.OrderStatusReceived:
		return pb.OrderStatus_RECEIVED
	case domain.OrderStatusPreparing:
		return pb.OrderStatus_PREPARING
	case domain.OrderStatusShipped:
		return pb.OrderStatus_SHIPPED
	case domain.OrderStatusDelivered:
		return pb.OrderStatus_DELIVERED
	case domain.OrderStatusCancelled:
		return pb.OrderStatus_CANCELLED
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

func domainCarrierToProto(carrier *domain.Carrier) pb.Carrier {
	if carrier == nil {
		return pb.Carrier_CARRIER_UNSPECIFIED
	}
	switch *carrier {
	case domain.CarrierYamato:
		return pb.Carrier_YAMATO
	case domain.CarrierSagawa:
		return pb.Carrier_SAGAWA
	case domain.CarrierJapanPost:
		return pb.Carrier_JAPAN_POST
	default:
		return pb.Carrier_CARRIER_UNSPECIFIED
	}
}

func protoCarrierToDomain(carrier pb.Carrier) *domain.Carrier {
	switch carrier {
	case pb.Carrier_YAMATO:
		c := domain.CarrierYamato
		return &c
	case pb.Carrier_SAGAWA:
		c := domain.CarrierSagawa
		return &c
	case pb.Carrier_JAPAN_POST:
		c := domain.CarrierJapanPost
		return &c
	default:
		return nil
	}
}
