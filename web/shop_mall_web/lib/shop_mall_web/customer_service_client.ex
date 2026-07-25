defmodule ShopMallWeb.CustomerServiceClient do
  @moduledoc """
  顧客サービス(customer service)への gRPC 呼び出しをまとめたクライアント。
  マイページ・カート・お気に入り・レビュー・商品ページから利用する。
  """

  alias CustomerService.V1.CustomerService.Stub
  alias CustomerService.V1, as: PB

  # ---- プロフィール ----

  def get_profile(customer_id) do
    call(fn ch -> Stub.get_profile(ch, %PB.GetProfileRequest{customer_id: customer_id}) end)
  end

  def update_profile(attrs) do
    request = %PB.UpdateProfileRequest{
      customer_id: attrs[:customer_id],
      first_name: attrs[:first_name] || "",
      last_name: attrs[:last_name] || "",
      phone: attrs[:phone] || "",
      birth_date: attrs[:birth_date] || "",
      gender: attrs[:gender] || :GENDER_UNSPECIFIED
    }

    call(fn ch -> Stub.update_profile(ch, request) end)
  end

  def upload_profile_image(customer_id, image_data) do
    request = %PB.UploadProfileImageRequest{customer_id: customer_id, image_data: image_data}
    call(fn ch -> Stub.upload_profile_image(ch, request) end)
  end

  # ---- 住所 ----

  def search_postal_code(postal_code) do
    call(fn ch ->
      Stub.search_postal_code(ch, %PB.SearchPostalCodeRequest{postal_code: postal_code})
    end)
  end

  def register_address(attrs) do
    request = %PB.RegisterAddressRequest{
      customer_id: attrs[:customer_id],
      address_name: attrs[:address_name] || "自宅",
      postal_code: attrs[:postal_code] || "",
      prefecture: attrs[:prefecture] || "",
      city: attrs[:city] || "",
      address_line1: attrs[:address_line1] || "",
      address_line2: attrs[:address_line2] || "",
      recipient_name: attrs[:recipient_name] || "",
      recipient_phone: attrs[:recipient_phone] || "",
      is_default: attrs[:is_default] || true
    }

    call(fn ch -> Stub.register_address(ch, request) end)
  end

  def update_address(attrs) do
    request = %PB.UpdateAddressRequest{
      address_id: attrs[:address_id],
      customer_id: attrs[:customer_id],
      address_name: attrs[:address_name] || "自宅",
      postal_code: attrs[:postal_code] || "",
      prefecture: attrs[:prefecture] || "",
      city: attrs[:city] || "",
      address_line1: attrs[:address_line1] || "",
      address_line2: attrs[:address_line2] || "",
      recipient_name: attrs[:recipient_name] || "",
      recipient_phone: attrs[:recipient_phone] || "",
      is_default: true
    }

    call(fn ch -> Stub.update_address(ch, request) end)
  end

  def delete_address(address_id, customer_id) do
    request = %PB.DeleteAddressRequest{address_id: address_id, customer_id: customer_id}
    call(fn ch -> Stub.delete_address(ch, request) end)
  end

  # ---- カート ----

  def add_to_cart(customer_id, product_id, quantity) do
    request = %PB.AddToCartRequest{
      customer_id: customer_id,
      product_id: product_id,
      quantity: quantity
    }

    call(fn ch -> Stub.add_to_cart(ch, request) end)
  end

  def get_cart(customer_id) do
    call(fn ch -> Stub.get_cart(ch, %PB.GetCartRequest{customer_id: customer_id}) end)
  end

  def update_cart_item_quantity(cart_item_id, customer_id, quantity) do
    request = %PB.UpdateCartItemQuantityRequest{
      cart_item_id: cart_item_id,
      customer_id: customer_id,
      quantity: quantity
    }

    call(fn ch -> Stub.update_cart_item_quantity(ch, request) end)
  end

  def remove_from_cart(cart_item_id, customer_id) do
    request = %PB.RemoveFromCartRequest{cart_item_id: cart_item_id, customer_id: customer_id}
    call(fn ch -> Stub.remove_from_cart(ch, request) end)
  end

  def merge_guest_cart(customer_id, session_id) do
    request = %PB.MergeGuestCartRequest{customer_id: customer_id, session_id: session_id}
    call(fn ch -> Stub.merge_guest_cart(ch, request) end)
  end

  # ---- お気に入り ----

  def add_to_favorite(customer_id, product_id) do
    request = %PB.AddToFavoriteRequest{
      customer_id: customer_id,
      product_id: product_id,
      notify_on_restock: true
    }

    call(fn ch -> Stub.add_to_favorite(ch, request) end)
  end

  def get_favorites(customer_id) do
    request = %PB.GetFavoritesRequest{
      customer_id: customer_id,
      sort_by: "created_at",
      sort_order: "desc"
    }

    call(fn ch -> Stub.get_favorites(ch, request) end)
  end

  def remove_from_favorite(favorite_id, customer_id) do
    request = %PB.RemoveFromFavoriteRequest{favorite_id: favorite_id, customer_id: customer_id}
    call(fn ch -> Stub.remove_from_favorite(ch, request) end)
  end

  # ---- 支払い方法 ----

  def register_payment_method(customer_id, stripe_payment_method_id) do
    request = %PB.RegisterPaymentMethodRequest{
      customer_id: customer_id,
      stripe_payment_method_id: stripe_payment_method_id,
      is_default: true
    }

    call(fn ch -> Stub.register_payment_method(ch, request) end)
  end

  def delete_payment_method(payment_method_id, customer_id) do
    request = %PB.DeletePaymentMethodRequest{
      payment_method_id: payment_method_id,
      customer_id: customer_id
    }

    call(fn ch -> Stub.delete_payment_method(ch, request) end)
  end

  # ---- レビュー ----

  def post_review(attrs) do
    request = %PB.PostReviewRequest{
      customer_id: attrs[:customer_id],
      product_id: attrs[:product_id],
      order_id: attrs[:order_id] || "",
      rating: attrs[:rating] || 5,
      review_text: attrs[:review_text] || ""
    }

    call(fn ch -> Stub.post_review(ch, request) end)
  end

  def update_review(attrs) do
    request = %PB.UpdateReviewRequest{
      review_id: attrs[:review_id],
      customer_id: attrs[:customer_id],
      rating: attrs[:rating] || 5,
      review_text: attrs[:review_text] || ""
    }

    call(fn ch -> Stub.update_review(ch, request) end)
  end

  def get_my_reviews(customer_id) do
    call(fn ch -> Stub.get_my_reviews(ch, %PB.GetMyReviewsRequest{customer_id: customer_id}) end)
  end

  # ---- 注文履歴(customer サービス経由で order へ委譲される) ----

  def get_order_history(customer_id) do
    request = %PB.GetOrderHistoryRequest{customer_id: customer_id, page: 1, page_size: 50}
    call(fn ch -> Stub.get_order_history(ch, request) end)
  end

  def get_order_detail(order_id, customer_id) do
    request = %PB.GetOrderDetailRequest{order_id: order_id, customer_id: customer_id}
    call(fn ch -> Stub.get_order_detail(ch, request) end)
  end

  def request_order_cancel(order_id, customer_id, reason) do
    request = %PB.RequestOrderCancelRequest{
      order_id: order_id,
      customer_id: customer_id,
      cancel_reason: reason
    }

    call(fn ch -> Stub.request_order_cancel(ch, request) end)
  end

  def reorder_from_history(order_id, customer_id) do
    request = %PB.ReorderFromHistoryRequest{order_id: order_id, customer_id: customer_id}
    call(fn ch -> Stub.reorder_from_history(ch, request) end)
  end

  # ---- 共通 ----

  defp call(fun) do
    host = System.get_env("CUSTOMER_SERVICE_HOST", "localhost")
    port = System.get_env("CUSTOMER_SERVICE_PORT", "50052")

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
        {:error, "顧客サービスに接続できません: #{inspect(reason)}"}
    end
  end
end
