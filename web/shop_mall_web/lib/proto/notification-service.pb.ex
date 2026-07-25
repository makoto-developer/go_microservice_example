# microservices/notification/proto から pb2ex.py で自動生成した最小スタブ。
defmodule NotificationService.V1.PreviewEmailTemplateRequest.VariablesEntry do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.PreviewEmailTemplateRequest.VariablesEntry",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:key, 1, type: :string)
  field(:value, 2, type: :string)
end

defmodule NotificationService.V1.SendBulkEmailRequest.VariablesEntry do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.SendBulkEmailRequest.VariablesEntry",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:key, 1, type: :string)
  field(:value, 2, type: :string)
end

defmodule NotificationService.V1.SendEmailRequest.VariablesEntry do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.SendEmailRequest.VariablesEntry",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:key, 1, type: :string)
  field(:value, 2, type: :string)
end

defmodule NotificationService.V1.SendPushNotificationRequest.VariablesEntry do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.SendPushNotificationRequest.VariablesEntry",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:key, 1, type: :string)
  field(:value, 2, type: :string)
end

defmodule NotificationService.V1.SendPushNotificationRequest.DataEntry do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.SendPushNotificationRequest.DataEntry",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:key, 1, type: :string)
  field(:value, 2, type: :string)
end

defmodule NotificationService.V1.NotificationChannel do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "notification_service.v1.NotificationChannel",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:NOTIFICATION_CHANNEL_UNSPECIFIED, 0)
  field(:EMAIL, 1)
  field(:PUSH, 2)
  field(:SMS, 3)
end

defmodule NotificationService.V1.NotificationFrequency do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "notification_service.v1.NotificationFrequency",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:NOTIFICATION_FREQUENCY_UNSPECIFIED, 0)
  field(:IMMEDIATE, 1)
  field(:DAILY_DIGEST, 2)
  field(:DISABLED, 3)
end

defmodule NotificationService.V1.NotificationStatus do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "notification_service.v1.NotificationStatus",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:NOTIFICATION_STATUS_UNSPECIFIED, 0)
  field(:PENDING, 1)
  field(:SENT, 2)
  field(:FAILED, 3)
  field(:CANCELLED, 4)
end

defmodule NotificationService.V1.NotificationType do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "notification_service.v1.NotificationType",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:NOTIFICATION_TYPE_UNSPECIFIED, 0)
  field(:USER_REGISTERED, 1)
  field(:EMAIL_VERIFIED, 2)
  field(:ORDER_CONFIRMED, 3)
  field(:ORDER_SHIPPED, 4)
  field(:ORDER_DELIVERED, 5)
  field(:ORDER_CANCELLED, 6)
  field(:PAYMENT_COMPLETED, 7)
  field(:PAYMENT_FAILED, 8)
  field(:REFUND_COMPLETED, 9)
  field(:SHOP_APPROVED, 10)
  field(:SHOP_REJECTED, 11)
  field(:STOCK_LOW_ALERT, 12)
  field(:STOCK_OUT_ALERT, 13)
  field(:STOCK_RESTORED, 14)
  field(:CHAT_MESSAGE_RECEIVED, 15)
  field(:PASSWORD_RESET, 16)
  field(:CAMPAIGN_STARTED, 17)
end

defmodule NotificationService.V1.Platform do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "notification_service.v1.Platform",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:PLATFORM_UNSPECIFIED, 0)
  field(:IOS, 1)
  field(:ANDROID, 2)
  field(:WEB, 3)
end

defmodule NotificationService.V1.CreateEmailTemplateRequest do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.CreateEmailTemplateRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:template_key, 1, type: :string, json_name: "templateKey")
  field(:subject_template, 2, type: :string, json_name: "subjectTemplate")
  field(:html_template, 3, type: :string, json_name: "htmlTemplate")
  field(:text_template, 4, type: :string, json_name: "textTemplate")
  field(:variables, 5, repeated: true, type: :string)
end

defmodule NotificationService.V1.CreateEmailTemplateResponse do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.CreateEmailTemplateResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule NotificationService.V1.GetNotificationHistoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.GetNotificationHistoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:channel, 2, type: NotificationService.V1.NotificationChannel, enum: true)
  field(:status, 3, type: NotificationService.V1.NotificationStatus, enum: true)
  field(:date_from, 4, type: :string, json_name: "dateFrom")
  field(:date_to, 5, type: :string, json_name: "dateTo")
  field(:page, 6, type: :int32)
  field(:page_size, 7, type: :int32, json_name: "pageSize")
end

defmodule NotificationService.V1.GetNotificationHistoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.GetNotificationHistoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule NotificationService.V1.GetNotificationPreferenceRequest do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.GetNotificationPreferenceRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
end

defmodule NotificationService.V1.GetNotificationPreferenceResponse do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.GetNotificationPreferenceResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule NotificationService.V1.PreviewEmailTemplateRequest do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.PreviewEmailTemplateRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:template_key, 1, type: :string, json_name: "templateKey")

  field(:variables, 2,
    repeated: true,
    type: NotificationService.V1.PreviewEmailTemplateRequest.VariablesEntry,
    map: true
  )
end

defmodule NotificationService.V1.PreviewEmailTemplateResponse do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.PreviewEmailTemplateResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule NotificationService.V1.RefreshDeviceTokenRequest do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.RefreshDeviceTokenRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:device_id, 1, type: :string, json_name: "deviceId")
  field(:new_token, 2, type: :string, json_name: "newToken")
end

defmodule NotificationService.V1.RefreshDeviceTokenResponse do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.RefreshDeviceTokenResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule NotificationService.V1.RegisterDeviceTokenRequest do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.RegisterDeviceTokenRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:device_id, 2, type: :string, json_name: "deviceId")
  field(:platform, 3, type: NotificationService.V1.Platform, enum: true)
  field(:token, 4, type: :string)
end

defmodule NotificationService.V1.RegisterDeviceTokenResponse do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.RegisterDeviceTokenResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule NotificationService.V1.ResendNotificationRequest do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.ResendNotificationRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:notification_id, 1, type: :string, json_name: "notificationId")
end

defmodule NotificationService.V1.ResendNotificationResponse do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.ResendNotificationResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule NotificationService.V1.SendBulkEmailRequest do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.SendBulkEmailRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_ids, 1, repeated: true, type: :string, json_name: "userIds")
  field(:template_key, 2, type: :string, json_name: "templateKey")

  field(:variables, 3,
    repeated: true,
    type: NotificationService.V1.SendBulkEmailRequest.VariablesEntry,
    map: true
  )
end

defmodule NotificationService.V1.SendBulkEmailResponse do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.SendBulkEmailResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule NotificationService.V1.SendEmailRequest do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.SendEmailRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:email, 2, type: :string)
  field(:template_key, 3, type: :string, json_name: "templateKey")

  field(:variables, 4,
    repeated: true,
    type: NotificationService.V1.SendEmailRequest.VariablesEntry,
    map: true
  )

  field(:notification_type, 5,
    type: NotificationService.V1.NotificationType,
    enum: true,
    json_name: "notificationType"
  )
end

defmodule NotificationService.V1.SendEmailResponse do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.SendEmailResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule NotificationService.V1.SendPushNotificationRequest do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.SendPushNotificationRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:template_key, 2, type: :string, json_name: "templateKey")

  field(:variables, 3,
    repeated: true,
    type: NotificationService.V1.SendPushNotificationRequest.VariablesEntry,
    map: true
  )

  field(:notification_type, 4,
    type: NotificationService.V1.NotificationType,
    enum: true,
    json_name: "notificationType"
  )

  field(:data, 5,
    repeated: true,
    type: NotificationService.V1.SendPushNotificationRequest.DataEntry,
    map: true
  )
end

defmodule NotificationService.V1.SendPushNotificationResponse do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.SendPushNotificationResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule NotificationService.V1.UnregisterDeviceTokenRequest do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.UnregisterDeviceTokenRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:device_id, 1, type: :string, json_name: "deviceId")
end

defmodule NotificationService.V1.UnregisterDeviceTokenResponse do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.UnregisterDeviceTokenResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule NotificationService.V1.UpdateEmailTemplateRequest do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.UpdateEmailTemplateRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:template_id, 1, type: :string, json_name: "templateId")
  field(:subject_template, 2, type: :string, json_name: "subjectTemplate")
  field(:html_template, 3, type: :string, json_name: "htmlTemplate")
  field(:text_template, 4, type: :string, json_name: "textTemplate")
  field(:variables, 5, repeated: true, type: :string)
end

defmodule NotificationService.V1.UpdateEmailTemplateResponse do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.UpdateEmailTemplateResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule NotificationService.V1.UpdateNotificationPreferenceRequest do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.UpdateNotificationPreferenceRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:email_enabled, 2, type: :bool, json_name: "emailEnabled")
  field(:push_enabled, 3, type: :bool, json_name: "pushEnabled")
  field(:email_order_updates, 4, type: :bool, json_name: "emailOrderUpdates")
  field(:email_shop_updates, 5, type: :bool, json_name: "emailShopUpdates")
  field(:email_chat_messages, 6, type: :bool, json_name: "emailChatMessages")
  field(:push_order_updates, 7, type: :bool, json_name: "pushOrderUpdates")
  field(:push_stock_restored, 8, type: :bool, json_name: "pushStockRestored")
  field(:push_campaigns, 9, type: :bool, json_name: "pushCampaigns")
  field(:push_chat_messages, 10, type: :bool, json_name: "pushChatMessages")
  field(:frequency, 11, type: NotificationService.V1.NotificationFrequency, enum: true)
  field(:quiet_hours_start, 12, type: Google.Protobuf.Timestamp, json_name: "quietHoursStart")
  field(:quiet_hours_end, 13, type: Google.Protobuf.Timestamp, json_name: "quietHoursEnd")
end

defmodule NotificationService.V1.UpdateNotificationPreferenceResponse do
  @moduledoc false

  use Protobuf,
    full_name: "notification_service.v1.UpdateNotificationPreferenceResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule NotificationService.V1.NotificationService.Service do
  @moduledoc false

  use GRPC.Service,
    name: "notification_service.v1.NotificationService",
    protoc_gen_elixir_version: "0.16.0"

  rpc(
    :SendEmail,
    NotificationService.V1.SendEmailRequest,
    NotificationService.V1.SendEmailResponse
  )

  rpc(
    :SendBulkEmail,
    NotificationService.V1.SendBulkEmailRequest,
    NotificationService.V1.SendBulkEmailResponse
  )

  rpc(
    :SendPushNotification,
    NotificationService.V1.SendPushNotificationRequest,
    NotificationService.V1.SendPushNotificationResponse
  )

  rpc(
    :RegisterDeviceToken,
    NotificationService.V1.RegisterDeviceTokenRequest,
    NotificationService.V1.RegisterDeviceTokenResponse
  )

  rpc(
    :UnregisterDeviceToken,
    NotificationService.V1.UnregisterDeviceTokenRequest,
    NotificationService.V1.UnregisterDeviceTokenResponse
  )

  rpc(
    :RefreshDeviceToken,
    NotificationService.V1.RefreshDeviceTokenRequest,
    NotificationService.V1.RefreshDeviceTokenResponse
  )

  rpc(
    :CreateEmailTemplate,
    NotificationService.V1.CreateEmailTemplateRequest,
    NotificationService.V1.CreateEmailTemplateResponse
  )

  rpc(
    :UpdateEmailTemplate,
    NotificationService.V1.UpdateEmailTemplateRequest,
    NotificationService.V1.UpdateEmailTemplateResponse
  )

  rpc(
    :PreviewEmailTemplate,
    NotificationService.V1.PreviewEmailTemplateRequest,
    NotificationService.V1.PreviewEmailTemplateResponse
  )

  rpc(
    :GetNotificationPreference,
    NotificationService.V1.GetNotificationPreferenceRequest,
    NotificationService.V1.GetNotificationPreferenceResponse
  )

  rpc(
    :UpdateNotificationPreference,
    NotificationService.V1.UpdateNotificationPreferenceRequest,
    NotificationService.V1.UpdateNotificationPreferenceResponse
  )

  rpc(
    :GetNotificationHistory,
    NotificationService.V1.GetNotificationHistoryRequest,
    NotificationService.V1.GetNotificationHistoryResponse
  )

  rpc(
    :ResendNotification,
    NotificationService.V1.ResendNotificationRequest,
    NotificationService.V1.ResendNotificationResponse
  )
end

defmodule NotificationService.V1.NotificationService.Stub do
  @moduledoc false

  use GRPC.Stub, service: NotificationService.V1.NotificationService.Service
end
