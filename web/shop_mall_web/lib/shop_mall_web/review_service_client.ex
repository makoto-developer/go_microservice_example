defmodule ShopMallWeb.ReviewServiceClient do
  @moduledoc """
  レビューサービス(review service)への gRPC 呼び出しをまとめたクライアント。
  商品ページ・マイページ・オーナー返信・管理者モデレーションから利用する。
  """

  alias ReviewService.V1.ReviewService.Stub
  alias ReviewService.V1, as: PB

  # ---- 顧客 ----

  def post_review(attrs) do
    request = %PB.PostReviewRequest{
      product_id: attrs[:product_id],
      order_id: attrs[:order_id] || "",
      customer_id: attrs[:customer_id],
      nickname: attrs[:nickname] || "匿名",
      rating: attrs[:rating] || 5,
      title: attrs[:title] || "",
      content: attrs[:content] || ""
    }

    call(fn ch -> Stub.post_review(ch, request) end)
  end

  def update_review(review_id, customer_id, rating, content) do
    request = %PB.UpdateReviewRequest{
      review_id: review_id,
      customer_id: customer_id,
      rating: rating,
      content: content
    }

    call(fn ch -> Stub.update_review(ch, request) end)
  end

  def delete_review(review_id, customer_id) do
    call(fn ch ->
      Stub.delete_review(ch, %PB.DeleteReviewRequest{
        review_id: review_id,
        customer_id: customer_id
      })
    end)
  end

  def get_reviews_by_product(product_id) do
    request = %PB.GetReviewsByProductRequest{
      product_id: product_id,
      sort_by: "newest",
      page: 1,
      page_size: 20
    }

    call(fn ch -> Stub.get_reviews_by_product(ch, request) end)
  end

  def get_product_rating(product_id) do
    call(fn ch ->
      Stub.get_product_rating(ch, %PB.GetProductRatingRequest{product_id: product_id})
    end)
  end

  def get_review_detail(review_id) do
    call(fn ch ->
      Stub.get_review_detail(ch, %PB.GetReviewDetailRequest{review_id: review_id})
    end)
  end

  def get_my_reviews(customer_id) do
    call(fn ch -> Stub.get_my_reviews(ch, %PB.GetMyReviewsRequest{customer_id: customer_id}) end)
  end

  def mark_helpful(review_id, user_id) do
    call(fn ch ->
      Stub.mark_review_helpful(ch, %PB.MarkReviewHelpfulRequest{
        review_id: review_id,
        user_id: user_id
      })
    end)
  end

  def unmark_helpful(review_id, user_id) do
    call(fn ch ->
      Stub.unmark_review_helpful(ch, %PB.UnmarkReviewHelpfulRequest{
        review_id: review_id,
        user_id: user_id
      })
    end)
  end

  def report_review(review_id, reporter_id, reason) do
    request = %PB.ReportReviewRequest{
      review_id: review_id,
      reporter_id: reporter_id,
      reason: reason,
      description: reason
    }

    call(fn ch -> Stub.report_review(ch, request) end)
  end

  # ---- 店舗オーナー(返信) ----

  def post_shop_reply(review_id, shop_id, content) do
    request = %PB.PostShopReplyRequest{review_id: review_id, shop_id: shop_id, content: content}
    call(fn ch -> Stub.post_shop_reply(ch, request) end)
  end

  def update_shop_reply(reply_id, shop_id, content) do
    request = %PB.UpdateShopReplyRequest{reply_id: reply_id, shop_id: shop_id, content: content}
    call(fn ch -> Stub.update_shop_reply(ch, request) end)
  end

  def delete_shop_reply(reply_id, shop_id) do
    call(fn ch ->
      Stub.delete_shop_reply(ch, %PB.DeleteShopReplyRequest{reply_id: reply_id, shop_id: shop_id})
    end)
  end

  # ---- 管理者(モデレーション) ----

  def get_pending_reviews do
    call(fn ch ->
      Stub.get_pending_reviews(ch, %PB.GetPendingReviewsRequest{page: 1, page_size: 20})
    end)
  end

  def approve_review(review_id, admin_id) do
    call(fn ch ->
      Stub.approve_review(ch, %PB.ApproveReviewRequest{review_id: review_id, admin_id: admin_id})
    end)
  end

  def reject_review(review_id, admin_id, reason) do
    request = %PB.RejectReviewRequest{review_id: review_id, admin_id: admin_id, reason: reason}
    call(fn ch -> Stub.reject_review(ch, request) end)
  end

  def delete_review_by_admin(review_id, admin_id, reason) do
    request = %PB.DeleteReviewByAdminRequest{
      review_id: review_id,
      admin_id: admin_id,
      reason: reason
    }

    call(fn ch -> Stub.delete_review_by_admin(ch, request) end)
  end

  def get_review_reports(admin_id) do
    request = %PB.GetReviewReportsRequest{
      admin_id: admin_id,
      status: "open",
      page: 1,
      page_size: 20
    }

    call(fn ch -> Stub.get_review_reports(ch, request) end)
  end

  def resolve_review_report(report_id, admin_id, action, note) do
    request = %PB.ResolveReviewReportRequest{
      report_id: report_id,
      admin_id: admin_id,
      action: action,
      note: note
    }

    call(fn ch -> Stub.resolve_review_report(ch, request) end)
  end

  defp call(fun) do
    host = System.get_env("REVIEW_SERVICE_HOST", "localhost")
    port = System.get_env("REVIEW_SERVICE_PORT", "20108")

    case GRPC.Stub.connect("#{host}:#{port}") do
      {:ok, channel} ->
        try do
          case fun.(channel) do
            {:ok, response} -> {:ok, response}
            {:error, %GRPC.RPCError{message: message}} -> {:error, message}
            {:error, reason} -> {:error, inspect(reason)}
          end
        after
          GRPC.Stub.disconnect(channel)
        end

      {:error, reason} ->
        {:error, "レビューサービスに接続できません: #{inspect(reason)}"}
    end
  end
end
