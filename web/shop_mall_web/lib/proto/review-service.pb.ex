# microservices/review/proto から pb2ex.py で自動生成した最小スタブ。
defmodule ReviewService.V1.ReportReason do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "review_service.v1.ReportReason",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:REPORT_REASON_UNSPECIFIED, 0)
  field(:INAPPROPRIATE_CONTENT, 1)
  field(:SPAM, 2)
  field(:OFF_TOPIC, 3)
  field(:FALSE_INFORMATION, 4)
  field(:PERSONAL_INFORMATION, 5)
  field(:OTHER, 6)
end

defmodule ReviewService.V1.ReportStatus do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "review_service.v1.ReportStatus",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:REPORT_STATUS_UNSPECIFIED, 0)
  field(:REPORT_STATUS_PENDING, 1)
  field(:REPORT_STATUS_REVIEWED, 2)
  field(:REPORT_STATUS_ACTION_TAKEN, 3)
  field(:REPORT_STATUS_DISMISSED, 4)
end

defmodule ReviewService.V1.ApproveReviewRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.ApproveReviewRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:review_id, 1, type: :string, json_name: "reviewId")
  field(:admin_id, 2, type: :string, json_name: "adminId")
end

defmodule ReviewService.V1.ApproveReviewResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.ApproveReviewResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.DeleteReviewByAdminRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.DeleteReviewByAdminRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:review_id, 1, type: :string, json_name: "reviewId")
  field(:admin_id, 2, type: :string, json_name: "adminId")
  field(:reason, 3, type: :string)
end

defmodule ReviewService.V1.DeleteReviewByAdminResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.DeleteReviewByAdminResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.DeleteReviewRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.DeleteReviewRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:review_id, 1, type: :string, json_name: "reviewId")
  field(:customer_id, 2, type: :string, json_name: "customerId")
end

defmodule ReviewService.V1.DeleteReviewResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.DeleteReviewResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.DeleteShopReplyRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.DeleteShopReplyRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:reply_id, 1, type: :string, json_name: "replyId")
  field(:shop_id, 2, type: :string, json_name: "shopId")
end

defmodule ReviewService.V1.DeleteShopReplyResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.DeleteShopReplyResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.GetMyReviewsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.GetMyReviewsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
end

defmodule ReviewService.V1.GetMyReviewsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.GetMyReviewsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.GetPendingReviewsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.GetPendingReviewsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:page, 1, type: :int32)
  field(:page_size, 2, type: :int32, json_name: "pageSize")
end

defmodule ReviewService.V1.GetPendingReviewsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.GetPendingReviewsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.GetProductRatingRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.GetProductRatingRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
end

defmodule ReviewService.V1.GetProductRatingResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.GetProductRatingResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.GetReviewDetailRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.GetReviewDetailRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:review_id, 1, type: :string, json_name: "reviewId")
end

defmodule ReviewService.V1.GetReviewDetailResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.GetReviewDetailResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.GetReviewReportsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.GetReviewReportsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:admin_id, 1, type: :string, json_name: "adminId")
  field(:status, 2, type: ReviewService.V1.ReportStatus, enum: true)
  field(:page, 3, type: :int32)
  field(:page_size, 4, type: :int32, json_name: "pageSize")
end

defmodule ReviewService.V1.GetReviewReportsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.GetReviewReportsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.GetReviewsByProductRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.GetReviewsByProductRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
  field(:sort_by, 2, type: :string, json_name: "sortBy")
  field(:page, 3, type: :int32)
  field(:page_size, 4, type: :int32, json_name: "pageSize")
end

defmodule ReviewService.V1.GetReviewsByProductResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.GetReviewsByProductResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.MarkReviewHelpfulRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.MarkReviewHelpfulRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:review_id, 1, type: :string, json_name: "reviewId")
  field(:user_id, 2, type: :string, json_name: "userId")
end

defmodule ReviewService.V1.MarkReviewHelpfulResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.MarkReviewHelpfulResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.PostReviewRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.PostReviewRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
  field(:order_id, 2, type: :string, json_name: "orderId")
  field(:customer_id, 3, type: :string, json_name: "customerId")
  field(:nickname, 4, type: :string)
  field(:rating, 5, type: :int32)
  field(:title, 6, type: :string)
  field(:content, 7, type: :string)
  field(:image_urls, 8, repeated: true, type: :string, json_name: "imageUrls")
end

defmodule ReviewService.V1.PostReviewResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.PostReviewResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.PostShopReplyRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.PostShopReplyRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:review_id, 1, type: :string, json_name: "reviewId")
  field(:shop_id, 2, type: :string, json_name: "shopId")
  field(:content, 3, type: :string)
end

defmodule ReviewService.V1.PostShopReplyResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.PostShopReplyResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.RejectReviewRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.RejectReviewRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:review_id, 1, type: :string, json_name: "reviewId")
  field(:admin_id, 2, type: :string, json_name: "adminId")
  field(:reason, 3, type: :string)
end

defmodule ReviewService.V1.RejectReviewResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.RejectReviewResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.ReportReviewRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.ReportReviewRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:review_id, 1, type: :string, json_name: "reviewId")
  field(:reporter_id, 2, type: :string, json_name: "reporterId")
  field(:reason, 3, type: ReviewService.V1.ReportReason, enum: true)
  field(:description, 4, type: :string)
end

defmodule ReviewService.V1.ReportReviewResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.ReportReviewResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.ResolveReviewReportRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.ResolveReviewReportRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:report_id, 1, type: :string, json_name: "reportId")
  field(:admin_id, 2, type: :string, json_name: "adminId")
  field(:action, 3, type: :string)
  field(:note, 4, type: :string)
end

defmodule ReviewService.V1.ResolveReviewReportResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.ResolveReviewReportResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.UnmarkReviewHelpfulRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.UnmarkReviewHelpfulRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:review_id, 1, type: :string, json_name: "reviewId")
  field(:user_id, 2, type: :string, json_name: "userId")
end

defmodule ReviewService.V1.UnmarkReviewHelpfulResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.UnmarkReviewHelpfulResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.UpdateReviewRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.UpdateReviewRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:review_id, 1, type: :string, json_name: "reviewId")
  field(:customer_id, 2, type: :string, json_name: "customerId")
  field(:rating, 3, type: :int32)
  field(:title, 4, type: :string)
  field(:content, 5, type: :string)
  field(:image_urls, 6, repeated: true, type: :string, json_name: "imageUrls")
end

defmodule ReviewService.V1.UpdateReviewResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.UpdateReviewResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.UpdateShopReplyRequest do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.UpdateShopReplyRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:reply_id, 1, type: :string, json_name: "replyId")
  field(:shop_id, 2, type: :string, json_name: "shopId")
  field(:content, 3, type: :string)
end

defmodule ReviewService.V1.UpdateShopReplyResponse do
  @moduledoc false

  use Protobuf,
    full_name: "review_service.v1.UpdateShopReplyResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ReviewService.V1.ReviewService.Service do
  @moduledoc false

  use GRPC.Service,
    name: "review_service.v1.ReviewService",
    protoc_gen_elixir_version: "0.16.0"

  rpc(:PostReview, ReviewService.V1.PostReviewRequest, ReviewService.V1.PostReviewResponse)
  rpc(:UpdateReview, ReviewService.V1.UpdateReviewRequest, ReviewService.V1.UpdateReviewResponse)
  rpc(:DeleteReview, ReviewService.V1.DeleteReviewRequest, ReviewService.V1.DeleteReviewResponse)

  rpc(
    :DeleteReviewByAdmin,
    ReviewService.V1.DeleteReviewByAdminRequest,
    ReviewService.V1.DeleteReviewByAdminResponse
  )

  rpc(
    :GetReviewsByProduct,
    ReviewService.V1.GetReviewsByProductRequest,
    ReviewService.V1.GetReviewsByProductResponse
  )

  rpc(
    :GetReviewDetail,
    ReviewService.V1.GetReviewDetailRequest,
    ReviewService.V1.GetReviewDetailResponse
  )

  rpc(:GetMyReviews, ReviewService.V1.GetMyReviewsRequest, ReviewService.V1.GetMyReviewsResponse)

  rpc(
    :GetProductRating,
    ReviewService.V1.GetProductRatingRequest,
    ReviewService.V1.GetProductRatingResponse
  )

  rpc(
    :GetPendingReviews,
    ReviewService.V1.GetPendingReviewsRequest,
    ReviewService.V1.GetPendingReviewsResponse
  )

  rpc(
    :ApproveReview,
    ReviewService.V1.ApproveReviewRequest,
    ReviewService.V1.ApproveReviewResponse
  )

  rpc(:RejectReview, ReviewService.V1.RejectReviewRequest, ReviewService.V1.RejectReviewResponse)

  rpc(
    :MarkReviewHelpful,
    ReviewService.V1.MarkReviewHelpfulRequest,
    ReviewService.V1.MarkReviewHelpfulResponse
  )

  rpc(
    :UnmarkReviewHelpful,
    ReviewService.V1.UnmarkReviewHelpfulRequest,
    ReviewService.V1.UnmarkReviewHelpfulResponse
  )

  rpc(
    :PostShopReply,
    ReviewService.V1.PostShopReplyRequest,
    ReviewService.V1.PostShopReplyResponse
  )

  rpc(
    :UpdateShopReply,
    ReviewService.V1.UpdateShopReplyRequest,
    ReviewService.V1.UpdateShopReplyResponse
  )

  rpc(
    :DeleteShopReply,
    ReviewService.V1.DeleteShopReplyRequest,
    ReviewService.V1.DeleteShopReplyResponse
  )

  rpc(:ReportReview, ReviewService.V1.ReportReviewRequest, ReviewService.V1.ReportReviewResponse)

  rpc(
    :GetReviewReports,
    ReviewService.V1.GetReviewReportsRequest,
    ReviewService.V1.GetReviewReportsResponse
  )

  rpc(
    :ResolveReviewReport,
    ReviewService.V1.ResolveReviewReportRequest,
    ReviewService.V1.ResolveReviewReportResponse
  )
end

defmodule ReviewService.V1.ReviewService.Stub do
  @moduledoc false

  use GRPC.Stub, service: ReviewService.V1.ReviewService.Service
end
