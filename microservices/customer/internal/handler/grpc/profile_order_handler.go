package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/client"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/customer/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WithReviewRepo は GetMyReviews 用のリポジトリを注入する。
func (h *CustomerServiceHandler) WithReviewRepo(repo repository.ReviewRepository) *CustomerServiceHandler {
	h.reviewRepo = repo
	return h
}

// WithOrderClient は注文履歴系 RPC の委譲先(order サービス)を注入する。
func (h *CustomerServiceHandler) WithOrderClient(oc client.OrderClient) *CustomerServiceHandler {
	h.orderClient = oc
	return h
}

func customerToProto(c *domain.Customer) *pb.Customer {
	out := &pb.Customer{
		Id:        c.ID.String(),
		UserId:    c.UserID.String(),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Phone:     c.PhoneNumber,
		BirthDate: formatDate(c.BirthDate),
		CreatedAt: timestampProto(c.CreatedAt),
		UpdatedAt: timestampProto(c.UpdatedAt),
	}
	if c.ProfileImageURL != nil {
		out.ProfileImageUrl = *c.ProfileImageURL
	}
	return out
}

// GetProfile は顧客プロフィールを返す。
func (h *CustomerServiceHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	customerID, err := uuid.Parse(req.GetCustomerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	output, err := h.getProfileUsecase.Execute(ctx, usecase.GetCustomerProfileInput{CustomerID: customerID})
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.GetProfileResponse{
		Customer: customerToProto(output.Customer),
		Message:  "Profile fetched",
	}, nil
}

// UpdateProfile は顧客プロフィール(氏名・電話・生年月日・性別)を更新する。
func (h *CustomerServiceHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	customerID, err := uuid.Parse(req.GetCustomerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}
	birthDate, err := parseDate(req.GetBirthDate())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birth_date (YYYY-MM-DD)")
	}

	output, err := h.updateProfileUsecase.Execute(ctx, usecase.UpdateCustomerProfileInput{
		CustomerID: customerID,
		FirstName:  req.GetFirstName(),
		LastName:   req.GetLastName(),
		Phone:      req.GetPhone(),
		BirthDate:  birthDate,
		Gender:     protoGenderToDomain(req.GetGender()),
	})
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.UpdateProfileResponse{
		Customer: customerToProto(output.Customer),
		Message:  "Profile updated",
	}, nil
}

// UploadProfileImage はプロフィール画像を保存する。
// このサンプルではストレージを持たないため、URL の発行のみ行う。
func (h *CustomerServiceHandler) UploadProfileImage(ctx context.Context, req *pb.UploadProfileImageRequest) (*pb.UploadProfileImageResponse, error) {
	customerID, err := uuid.Parse(req.GetCustomerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}
	if len(req.GetImageData()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "image_data is required")
	}

	base := fmt.Sprintf("/uploads/profiles/%s", customerID)
	return &pb.UploadProfileImageResponse{
		ImageUrl:         base + ".jpg",
		Thumbnail_100Url: base + "_100.jpg",
		Thumbnail_200Url: base + "_200.jpg",
		Message:          "Profile image uploaded",
	}, nil
}

// サンプル用の郵便番号辞書。実運用では郵便番号 API を呼ぶ想定。
var postalCodeTable = map[string][3]string{
	"1000001": {"東京都", "千代田区", "千代田"},
	"1500041": {"東京都", "渋谷区", "神南"},
	"5300001": {"大阪府", "大阪市北区", "梅田"},
	"4600002": {"愛知県", "名古屋市中区", "丸の内"},
}

// SearchPostalCode は郵便番号から住所を引く(住所入力フォームの自動補完用)。
func (h *CustomerServiceHandler) SearchPostalCode(ctx context.Context, req *pb.SearchPostalCodeRequest) (*pb.SearchPostalCodeResponse, error) {
	code := strings.ReplaceAll(req.GetPostalCode(), "-", "")
	if len(code) != 7 {
		return nil, status.Error(codes.InvalidArgument, "postal_code must be 7 digits")
	}
	if entry, ok := postalCodeTable[code]; ok {
		return &pb.SearchPostalCodeResponse{
			Prefecture:   entry[0],
			City:         entry[1],
			AddressLine1: entry[2],
		}, nil
	}
	return nil, status.Error(codes.NotFound, "postal code not found")
}

// GetOrderHistory は注文履歴を返す(order サービスへ委譲)。
func (h *CustomerServiceHandler) GetOrderHistory(ctx context.Context, req *pb.GetOrderHistoryRequest) (*pb.GetOrderHistoryResponse, error) {
	if h.orderClient == nil {
		return nil, status.Error(codes.Unavailable, "order service is not configured")
	}
	orders, err := h.orderClient.ListOrders(ctx, req.GetCustomerId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch order history: %v", err)
	}

	out := make([]*pb.OrderSummary, 0, len(orders))
	for _, o := range orders {
		amount, _ := parseAmount(o.GetTotalAmount())
		out = append(out, &pb.OrderSummary{
			OrderId:     o.GetId(),
			OrderNumber: o.GetOrderNumber(),
			Status:      o.GetStatus().String(),
			TotalAmount: amount,
			OrderedAt:   o.GetCreatedAt(),
		})
	}
	return &pb.GetOrderHistoryResponse{
		Orders:     out,
		TotalCount: int32(len(out)),
		Page:       1,
		PageSize:   int32(len(out)),
	}, nil
}

// GetOrderDetail は注文の詳細を返す(order サービスへ委譲)。
func (h *CustomerServiceHandler) GetOrderDetail(ctx context.Context, req *pb.GetOrderDetailRequest) (*pb.GetOrderDetailResponse, error) {
	if h.orderClient == nil {
		return nil, status.Error(codes.Unavailable, "order service is not configured")
	}
	order, err := h.orderClient.GetOrderDetail(ctx, req.GetOrderId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch order detail: %v", err)
	}
	amount, _ := parseAmount(order.GetTotalAmount())
	return &pb.GetOrderDetailResponse{
		Order: &pb.OrderSummary{
			OrderId:     order.GetId(),
			OrderNumber: order.GetOrderNumber(),
			Status:      order.GetStatus().String(),
			TotalAmount: amount,
			OrderedAt:   order.GetCreatedAt(),
		},
	}, nil
}

// RequestOrderCancel は注文キャンセルを受け付ける(order サービスへ委譲。決済済みなら返金される)。
func (h *CustomerServiceHandler) RequestOrderCancel(ctx context.Context, req *pb.RequestOrderCancelRequest) (*pb.RequestOrderCancelResponse, error) {
	if h.orderClient == nil {
		return nil, status.Error(codes.Unavailable, "order service is not configured")
	}
	if err := h.orderClient.CancelOrder(ctx, req.GetOrderId(), req.GetCustomerId(), req.GetCancelReason()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to cancel order: %v", err)
	}
	return &pb.RequestOrderCancelResponse{
		Success: true,
		Message: "Order cancelled. Refund will be processed if already paid.",
	}, nil
}

// ReorderFromHistory は過去の注文を再注文する(order サービスへ委譲)。
func (h *CustomerServiceHandler) ReorderFromHistory(ctx context.Context, req *pb.ReorderFromHistoryRequest) (*pb.ReorderFromHistoryResponse, error) {
	if h.orderClient == nil {
		return nil, status.Error(codes.Unavailable, "order service is not configured")
	}
	orderID, err := h.orderClient.Reorder(ctx, req.GetOrderId(), req.GetCustomerId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reorder: %v", err)
	}
	return &pb.ReorderFromHistoryResponse{
		Message: "Reorder created: " + orderID,
	}, nil
}

// GetMyReviews は自分が書いたレビューの一覧を返す。
func (h *CustomerServiceHandler) GetMyReviews(ctx context.Context, req *pb.GetMyReviewsRequest) (*pb.GetMyReviewsResponse, error) {
	if h.reviewRepo == nil {
		return nil, status.Error(codes.Unavailable, "review repository is not configured")
	}
	customerID, err := uuid.Parse(req.GetCustomerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}
	reviews, err := h.reviewRepo.ListByCustomer(ctx, customerID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list reviews: %v", err)
	}

	out := make([]*pb.Review, 0, len(reviews))
	for _, r := range reviews {
		out = append(out, &pb.Review{
			Id:            r.ID.String(),
			CustomerId:    r.CustomerID.String(),
			ProductId:     r.ProductID.String(),
			OrderId:       r.OrderID.String(),
			Rating:        int32(r.Rating),
			ReviewText:    r.ReviewText,
			EditableUntil: timestampProto(r.EditableUntil),
			CreatedAt:     timestampProto(r.CreatedAt),
			UpdatedAt:     timestampProto(r.UpdatedAt),
		})
	}
	return &pb.GetMyReviewsResponse{Reviews: out, TotalCount: int32(len(out))}, nil
}

func parseAmount(s string) (int64, error) {
	var n int64
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
