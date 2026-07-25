# microservices/customer/proto から pb2ex.py で自動生成した最小スタブ。
# フィールド番号・型は customer-service.pb.go を正とすること。
defmodule CustomerService.V1.Gender do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "customer_service.v1.Gender",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:GENDER_UNSPECIFIED, 0)
  field(:MALE, 1)
  field(:FEMALE, 2)
  field(:OTHER, 3)
  field(:PREFER_NOT_TO_SAY, 4)
end

defmodule CustomerService.V1.AddToCartRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.AddToCartRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
  field(:product_id, 2, type: :string, json_name: "productId")
  field(:variation_id, 3, type: :string, json_name: "variationId")
  field(:quantity, 4, type: :int32)
end

defmodule CustomerService.V1.AddToCartResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.AddToCartResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:cart_item, 1, type: CustomerService.V1.CartItem, json_name: "cartItem")
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.AddToFavoriteRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.AddToFavoriteRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
  field(:product_id, 2, type: :string, json_name: "productId")
  field(:notify_on_restock, 3, type: :bool, json_name: "notifyOnRestock")
end

defmodule CustomerService.V1.AddToFavoriteResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.AddToFavoriteResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:favorite, 1, type: CustomerService.V1.Favorite)
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.Address do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.Address",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:customer_id, 2, type: :string, json_name: "customerId")
  field(:address_name, 3, type: :string, json_name: "addressName")
  field(:postal_code, 4, type: :string, json_name: "postalCode")
  field(:prefecture, 5, type: :string)
  field(:city, 6, type: :string)
  field(:address_line1, 7, type: :string, json_name: "addressLine1")
  field(:address_line2, 8, type: :string, json_name: "addressLine2")
  field(:recipient_name, 9, type: :string, json_name: "recipientName")
  field(:recipient_phone, 10, type: :string, json_name: "recipientPhone")
  field(:is_default, 11, type: :bool, json_name: "isDefault")
  field(:deleted, 12, type: :bool)
  field(:created_at, 13, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 14, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule CustomerService.V1.CartItem do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.CartItem",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:customer_id, 2, type: :string, json_name: "customerId")
  field(:product_id, 3, type: :string, json_name: "productId")
  field(:variation_id, 4, type: :string, json_name: "variationId")
  field(:quantity, 5, type: :int32)
  field(:expires_at, 6, type: Google.Protobuf.Timestamp, json_name: "expiresAt")
  field(:created_at, 7, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 8, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule CustomerService.V1.Customer do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.Customer",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:user_id, 2, type: :string, json_name: "userId")
  field(:first_name, 3, type: :string, json_name: "firstName")
  field(:last_name, 4, type: :string, json_name: "lastName")
  field(:phone, 5, type: :string)
  field(:birth_date, 6, type: :string, json_name: "birthDate")
  field(:gender, 7, type: CustomerService.V1.Gender, enum: true)
  field(:profile_image_url, 8, type: :string, json_name: "profileImageUrl")
  field(:profile_thumbnail_100_url, 9, type: :string, json_name: "profileThumbnail100Url")
  field(:profile_thumbnail_200_url, 10, type: :string, json_name: "profileThumbnail200Url")
  field(:created_at, 11, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 12, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule CustomerService.V1.DeleteAddressRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.DeleteAddressRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:address_id, 1, type: :string, json_name: "addressId")
  field(:customer_id, 2, type: :string, json_name: "customerId")
end

defmodule CustomerService.V1.DeleteAddressResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.DeleteAddressResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.DeletePaymentMethodRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.DeletePaymentMethodRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:payment_method_id, 1, type: :string, json_name: "paymentMethodId")
  field(:customer_id, 2, type: :string, json_name: "customerId")
end

defmodule CustomerService.V1.DeletePaymentMethodResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.DeletePaymentMethodResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.Favorite do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.Favorite",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:customer_id, 2, type: :string, json_name: "customerId")
  field(:product_id, 3, type: :string, json_name: "productId")
  field(:notify_on_restock, 4, type: :bool, json_name: "notifyOnRestock")
  field(:created_at, 5, type: Google.Protobuf.Timestamp, json_name: "createdAt")
end

defmodule CustomerService.V1.GetCartRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.GetCartRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
end

defmodule CustomerService.V1.GetCartResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.GetCartResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:cart_items, 1, repeated: true, type: CustomerService.V1.CartItem, json_name: "cartItems")
  field(:total_quantity, 2, type: :int32, json_name: "totalQuantity")
end

defmodule CustomerService.V1.GetFavoritesRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.GetFavoritesRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
  field(:sort_by, 2, type: :string, json_name: "sortBy")
  field(:sort_order, 3, type: :string, json_name: "sortOrder")
end

defmodule CustomerService.V1.GetFavoritesResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.GetFavoritesResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:favorites, 1, repeated: true, type: CustomerService.V1.Favorite)
  field(:total_count, 2, type: :int32, json_name: "totalCount")
end

defmodule CustomerService.V1.GetMyReviewsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.GetMyReviewsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
end

defmodule CustomerService.V1.GetMyReviewsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.GetMyReviewsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:reviews, 1, repeated: true, type: CustomerService.V1.Review)
  field(:total_count, 2, type: :int32, json_name: "totalCount")
end

defmodule CustomerService.V1.GetOrderDetailRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.GetOrderDetailRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order_id, 1, type: :string, json_name: "orderId")
  field(:customer_id, 2, type: :string, json_name: "customerId")
end

defmodule CustomerService.V1.GetOrderDetailResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.GetOrderDetailResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order, 1, type: CustomerService.V1.OrderSummary)
  field(:item_details, 2, repeated: true, type: :string, json_name: "itemDetails")
end

defmodule CustomerService.V1.GetOrderHistoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.GetOrderHistoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
  field(:status_filter, 2, type: :string, json_name: "statusFilter")
  field(:date_from, 3, type: :string, json_name: "dateFrom")
  field(:date_to, 4, type: :string, json_name: "dateTo")
  field(:page, 5, type: :int32)
  field(:page_size, 6, type: :int32, json_name: "pageSize")
end

defmodule CustomerService.V1.GetOrderHistoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.GetOrderHistoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:orders, 1, repeated: true, type: CustomerService.V1.OrderSummary)
  field(:total_count, 2, type: :int32, json_name: "totalCount")
  field(:page, 3, type: :int32)
  field(:page_size, 4, type: :int32, json_name: "pageSize")
end

defmodule CustomerService.V1.GetProfileRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.GetProfileRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
end

defmodule CustomerService.V1.GetProfileResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.GetProfileResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer, 1, type: CustomerService.V1.Customer)
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.MergeGuestCartRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.MergeGuestCartRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
  field(:session_id, 2, type: :string, json_name: "sessionId")
end

defmodule CustomerService.V1.MergeGuestCartResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.MergeGuestCartResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:cart_items, 1, repeated: true, type: CustomerService.V1.CartItem, json_name: "cartItems")
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.OrderSummary do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.OrderSummary",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order_id, 1, type: :string, json_name: "orderId")
  field(:order_number, 2, type: :string, json_name: "orderNumber")
  field(:status, 3, type: :string)
  field(:total_amount, 4, type: :int64, json_name: "totalAmount")
  field(:ordered_at, 5, type: Google.Protobuf.Timestamp, json_name: "orderedAt")
end

defmodule CustomerService.V1.PaymentMethod do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.PaymentMethod",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:customer_id, 2, type: :string, json_name: "customerId")
  field(:stripe_payment_method_id, 3, type: :string, json_name: "stripePaymentMethodId")
  field(:card_last4, 4, type: :string, json_name: "cardLast4")
  field(:card_brand, 5, type: :string, json_name: "cardBrand")
  field(:card_exp_month, 6, type: :int32, json_name: "cardExpMonth")
  field(:card_exp_year, 7, type: :int32, json_name: "cardExpYear")
  field(:cardholder_name, 8, type: :string, json_name: "cardholderName")
  field(:is_default, 9, type: :bool, json_name: "isDefault")
  field(:created_at, 10, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 11, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule CustomerService.V1.PostReviewRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.PostReviewRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
  field(:product_id, 2, type: :string, json_name: "productId")
  field(:order_id, 3, type: :string, json_name: "orderId")
  field(:rating, 4, type: :int32)
  field(:review_text, 5, type: :string, json_name: "reviewText")
  field(:image_urls, 6, repeated: true, type: :string, json_name: "imageUrls")
end

defmodule CustomerService.V1.PostReviewResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.PostReviewResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:review, 1, type: CustomerService.V1.Review)
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.RegisterAddressRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.RegisterAddressRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
  field(:address_name, 2, type: :string, json_name: "addressName")
  field(:postal_code, 3, type: :string, json_name: "postalCode")
  field(:prefecture, 4, type: :string)
  field(:city, 5, type: :string)
  field(:address_line1, 6, type: :string, json_name: "addressLine1")
  field(:address_line2, 7, type: :string, json_name: "addressLine2")
  field(:recipient_name, 8, type: :string, json_name: "recipientName")
  field(:recipient_phone, 9, type: :string, json_name: "recipientPhone")
  field(:is_default, 10, type: :bool, json_name: "isDefault")
end

defmodule CustomerService.V1.RegisterAddressResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.RegisterAddressResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:address, 1, type: CustomerService.V1.Address)
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.RegisterPaymentMethodRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.RegisterPaymentMethodRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
  field(:stripe_payment_method_id, 2, type: :string, json_name: "stripePaymentMethodId")
  field(:is_default, 3, type: :bool, json_name: "isDefault")
end

defmodule CustomerService.V1.RegisterPaymentMethodResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.RegisterPaymentMethodResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:payment_method, 1, type: CustomerService.V1.PaymentMethod, json_name: "paymentMethod")
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.RemoveFromCartRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.RemoveFromCartRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:cart_item_id, 1, type: :string, json_name: "cartItemId")
  field(:customer_id, 2, type: :string, json_name: "customerId")
end

defmodule CustomerService.V1.RemoveFromCartResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.RemoveFromCartResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.RemoveFromFavoriteRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.RemoveFromFavoriteRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:favorite_id, 1, type: :string, json_name: "favoriteId")
  field(:customer_id, 2, type: :string, json_name: "customerId")
end

defmodule CustomerService.V1.RemoveFromFavoriteResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.RemoveFromFavoriteResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.ReorderFromHistoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.ReorderFromHistoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order_id, 1, type: :string, json_name: "orderId")
  field(:customer_id, 2, type: :string, json_name: "customerId")
end

defmodule CustomerService.V1.ReorderFromHistoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.ReorderFromHistoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:cart_items, 1, repeated: true, type: CustomerService.V1.CartItem, json_name: "cartItems")
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.RequestOrderCancelRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.RequestOrderCancelRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order_id, 1, type: :string, json_name: "orderId")
  field(:customer_id, 2, type: :string, json_name: "customerId")
  field(:cancel_reason, 3, type: :string, json_name: "cancelReason")
end

defmodule CustomerService.V1.RequestOrderCancelResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.RequestOrderCancelResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.Review do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.Review",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:customer_id, 2, type: :string, json_name: "customerId")
  field(:product_id, 3, type: :string, json_name: "productId")
  field(:order_id, 4, type: :string, json_name: "orderId")
  field(:rating, 5, type: :int32)
  field(:review_text, 6, type: :string, json_name: "reviewText")
  field(:image_urls, 7, repeated: true, type: :string, json_name: "imageUrls")
  field(:editable_until, 8, type: Google.Protobuf.Timestamp, json_name: "editableUntil")
  field(:created_at, 9, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 10, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule CustomerService.V1.SearchPostalCodeRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.SearchPostalCodeRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:postal_code, 1, type: :string, json_name: "postalCode")
end

defmodule CustomerService.V1.SearchPostalCodeResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.SearchPostalCodeResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:prefecture, 1, type: :string)
  field(:city, 2, type: :string)
  field(:address_line1, 3, type: :string, json_name: "addressLine1")
end

defmodule CustomerService.V1.UpdateAddressRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.UpdateAddressRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:address_id, 1, type: :string, json_name: "addressId")
  field(:customer_id, 2, type: :string, json_name: "customerId")
  field(:address_name, 3, type: :string, json_name: "addressName")
  field(:postal_code, 4, type: :string, json_name: "postalCode")
  field(:prefecture, 5, type: :string)
  field(:city, 6, type: :string)
  field(:address_line1, 7, type: :string, json_name: "addressLine1")
  field(:address_line2, 8, type: :string, json_name: "addressLine2")
  field(:recipient_name, 9, type: :string, json_name: "recipientName")
  field(:recipient_phone, 10, type: :string, json_name: "recipientPhone")
  field(:is_default, 11, type: :bool, json_name: "isDefault")
end

defmodule CustomerService.V1.UpdateAddressResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.UpdateAddressResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:address, 1, type: CustomerService.V1.Address)
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.UpdateCartItemQuantityRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.UpdateCartItemQuantityRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:cart_item_id, 1, type: :string, json_name: "cartItemId")
  field(:customer_id, 2, type: :string, json_name: "customerId")
  field(:quantity, 3, type: :int32)
end

defmodule CustomerService.V1.UpdateCartItemQuantityResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.UpdateCartItemQuantityResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:cart_item, 1, type: CustomerService.V1.CartItem, json_name: "cartItem")
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.UpdateProfileRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.UpdateProfileRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
  field(:first_name, 2, type: :string, json_name: "firstName")
  field(:last_name, 3, type: :string, json_name: "lastName")
  field(:phone, 4, type: :string)
  field(:birth_date, 5, type: :string, json_name: "birthDate")
  field(:gender, 6, type: CustomerService.V1.Gender, enum: true)
end

defmodule CustomerService.V1.UpdateProfileResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.UpdateProfileResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer, 1, type: CustomerService.V1.Customer)
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.UpdateReviewRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.UpdateReviewRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:review_id, 1, type: :string, json_name: "reviewId")
  field(:customer_id, 2, type: :string, json_name: "customerId")
  field(:rating, 3, type: :int32)
  field(:review_text, 4, type: :string, json_name: "reviewText")
  field(:image_urls, 5, repeated: true, type: :string, json_name: "imageUrls")
end

defmodule CustomerService.V1.UpdateReviewResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.UpdateReviewResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:review, 1, type: CustomerService.V1.Review)
  field(:message, 2, type: :string)
end

defmodule CustomerService.V1.UploadProfileImageRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.UploadProfileImageRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
  field(:image_data, 2, type: :bytes, json_name: "imageData")
end

defmodule CustomerService.V1.UploadProfileImageResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_service.v1.UploadProfileImageResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:image_url, 1, type: :string, json_name: "imageUrl")
  field(:thumbnail_100_url, 2, type: :string, json_name: "thumbnail100Url")
  field(:thumbnail_200_url, 3, type: :string, json_name: "thumbnail200Url")
  field(:message, 4, type: :string)
end

defmodule CustomerService.V1.CustomerService.Service do
  @moduledoc false

  use GRPC.Service,
    name: "customer_service.v1.CustomerService",
    protoc_gen_elixir_version: "0.16.0"

  rpc(:GetProfile, CustomerService.V1.GetProfileRequest, CustomerService.V1.GetProfileResponse)

  rpc(
    :UpdateProfile,
    CustomerService.V1.UpdateProfileRequest,
    CustomerService.V1.UpdateProfileResponse
  )

  rpc(
    :UploadProfileImage,
    CustomerService.V1.UploadProfileImageRequest,
    CustomerService.V1.UploadProfileImageResponse
  )

  rpc(
    :RegisterAddress,
    CustomerService.V1.RegisterAddressRequest,
    CustomerService.V1.RegisterAddressResponse
  )

  rpc(
    :UpdateAddress,
    CustomerService.V1.UpdateAddressRequest,
    CustomerService.V1.UpdateAddressResponse
  )

  rpc(
    :DeleteAddress,
    CustomerService.V1.DeleteAddressRequest,
    CustomerService.V1.DeleteAddressResponse
  )

  rpc(
    :SearchPostalCode,
    CustomerService.V1.SearchPostalCodeRequest,
    CustomerService.V1.SearchPostalCodeResponse
  )

  rpc(:AddToCart, CustomerService.V1.AddToCartRequest, CustomerService.V1.AddToCartResponse)
  rpc(:GetCart, CustomerService.V1.GetCartRequest, CustomerService.V1.GetCartResponse)

  rpc(
    :UpdateCartItemQuantity,
    CustomerService.V1.UpdateCartItemQuantityRequest,
    CustomerService.V1.UpdateCartItemQuantityResponse
  )

  rpc(
    :RemoveFromCart,
    CustomerService.V1.RemoveFromCartRequest,
    CustomerService.V1.RemoveFromCartResponse
  )

  rpc(
    :MergeGuestCart,
    CustomerService.V1.MergeGuestCartRequest,
    CustomerService.V1.MergeGuestCartResponse
  )

  rpc(
    :AddToFavorite,
    CustomerService.V1.AddToFavoriteRequest,
    CustomerService.V1.AddToFavoriteResponse
  )

  rpc(
    :GetFavorites,
    CustomerService.V1.GetFavoritesRequest,
    CustomerService.V1.GetFavoritesResponse
  )

  rpc(
    :RemoveFromFavorite,
    CustomerService.V1.RemoveFromFavoriteRequest,
    CustomerService.V1.RemoveFromFavoriteResponse
  )

  rpc(
    :GetOrderHistory,
    CustomerService.V1.GetOrderHistoryRequest,
    CustomerService.V1.GetOrderHistoryResponse
  )

  rpc(
    :GetOrderDetail,
    CustomerService.V1.GetOrderDetailRequest,
    CustomerService.V1.GetOrderDetailResponse
  )

  rpc(
    :RequestOrderCancel,
    CustomerService.V1.RequestOrderCancelRequest,
    CustomerService.V1.RequestOrderCancelResponse
  )

  rpc(
    :ReorderFromHistory,
    CustomerService.V1.ReorderFromHistoryRequest,
    CustomerService.V1.ReorderFromHistoryResponse
  )

  rpc(
    :RegisterPaymentMethod,
    CustomerService.V1.RegisterPaymentMethodRequest,
    CustomerService.V1.RegisterPaymentMethodResponse
  )

  rpc(
    :DeletePaymentMethod,
    CustomerService.V1.DeletePaymentMethodRequest,
    CustomerService.V1.DeletePaymentMethodResponse
  )

  rpc(:PostReview, CustomerService.V1.PostReviewRequest, CustomerService.V1.PostReviewResponse)

  rpc(
    :UpdateReview,
    CustomerService.V1.UpdateReviewRequest,
    CustomerService.V1.UpdateReviewResponse
  )

  rpc(
    :GetMyReviews,
    CustomerService.V1.GetMyReviewsRequest,
    CustomerService.V1.GetMyReviewsResponse
  )
end

defmodule CustomerService.V1.CustomerService.Stub do
  @moduledoc false

  use GRPC.Stub, service: CustomerService.V1.CustomerService.Service
end
