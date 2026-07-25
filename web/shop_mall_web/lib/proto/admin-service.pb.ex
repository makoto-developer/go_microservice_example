# microservices/admin/proto から pb2ex.py で自動生成した最小スタブ。
defmodule AdminService.V1.HealthStatus do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "admin_service.v1.HealthStatus",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:HEALTH_STATUS_UNSPECIFIED, 0)
  field(:HEALTHY, 1)
  field(:DEGRADED, 2)
  field(:UNHEALTHY, 3)
end

defmodule AdminService.V1.OperationType do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "admin_service.v1.OperationType",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:OPERATION_TYPE_UNSPECIFIED, 0)
  field(:USER_CREATED, 1)
  field(:USER_UPDATED, 2)
  field(:USER_ROLE_CHANGED, 3)
  field(:USER_SUSPENDED, 4)
  field(:USER_ACTIVATED, 5)
  field(:SHOP_APPROVED, 6)
  field(:SHOP_REJECTED, 7)
  field(:SHOP_SUSPENDED, 8)
  field(:SHOP_ACTIVATED, 9)
  field(:SETTING_UPDATED, 10)
  field(:CATEGORY_CREATED, 11)
  field(:CATEGORY_UPDATED, 12)
  field(:CATEGORY_DELETED, 13)
  field(:REVIEW_APPROVED, 14)
  field(:REVIEW_REJECTED, 15)
  field(:REVIEW_DELETED, 16)
end

defmodule AdminService.V1.SettingType do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "admin_service.v1.SettingType",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:SETTING_TYPE_UNSPECIFIED, 0)
  field(:STRING, 1)
  field(:NUMBER, 2)
  field(:BOOLEAN, 3)
  field(:JSON, 4)
end

defmodule AdminService.V1.Severity do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "admin_service.v1.Severity",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:SEVERITY_UNSPECIFIED, 0)
  field(:LOW, 1)
  field(:MEDIUM, 2)
  field(:HIGH, 3)
  field(:CRITICAL, 4)
end

defmodule AdminService.V1.WordContext do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "admin_service.v1.WordContext",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:WORD_CONTEXT_UNSPECIFIED, 0)
  field(:REVIEW, 1)
  field(:CHAT, 2)
  field(:ALL, 3)
end

defmodule AdminService.V1.ActivateShopRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.ActivateShopRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:admin_id, 2, type: :string, json_name: "adminId")
  field(:reason, 3, type: :string)
end

defmodule AdminService.V1.ActivateShopResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.ActivateShopResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.ActivateUserRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.ActivateUserRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:admin_id, 2, type: :string, json_name: "adminId")
  field(:reason, 3, type: :string)
end

defmodule AdminService.V1.ActivateUserResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.ActivateUserResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.AddForbiddenWordRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.AddForbiddenWordRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:admin_id, 1, type: :string, json_name: "adminId")
  field(:word, 2, type: :string)
  field(:context, 3, type: AdminService.V1.WordContext, enum: true)
  field(:severity, 4, type: AdminService.V1.Severity, enum: true)
end

defmodule AdminService.V1.AddForbiddenWordResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.AddForbiddenWordResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.ApproveShopRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.ApproveShopRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:admin_id, 2, type: :string, json_name: "adminId")
end

defmodule AdminService.V1.ApproveShopResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.ApproveShopResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.Category do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.Category",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:name, 2, type: :string)
  field(:parent_id, 3, type: :string, json_name: "parentId")
  field(:level, 4, type: :int32)
  field(:display_order, 5, type: :int32, json_name: "displayOrder")
  field(:is_active, 6, type: :bool, json_name: "isActive")
  field(:created_at, 7, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 8, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule AdminService.V1.ChangeUserRoleRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.ChangeUserRoleRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:admin_id, 2, type: :string, json_name: "adminId")
  field(:new_role, 3, type: :string, json_name: "newRole")
  field(:reason, 4, type: :string)
end

defmodule AdminService.V1.ChangeUserRoleResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.ChangeUserRoleResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.CreateCategoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.CreateCategoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:admin_id, 1, type: :string, json_name: "adminId")
  field(:name, 2, type: :string)
  field(:parent_id, 3, type: :string, json_name: "parentId")
  field(:display_order, 4, type: :int32, json_name: "displayOrder")
end

defmodule AdminService.V1.CreateCategoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.CreateCategoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.DeleteCategoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.DeleteCategoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:admin_id, 1, type: :string, json_name: "adminId")
  field(:category_id, 2, type: :string, json_name: "categoryId")
end

defmodule AdminService.V1.DeleteCategoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.DeleteCategoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.DeleteForbiddenWordRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.DeleteForbiddenWordRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:admin_id, 1, type: :string, json_name: "adminId")
  field(:word_id, 2, type: :string, json_name: "wordId")
end

defmodule AdminService.V1.DeleteForbiddenWordResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.DeleteForbiddenWordResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.ExportAuditLogsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.ExportAuditLogsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:date_from, 1, type: :string, json_name: "dateFrom")
  field(:date_to, 2, type: :string, json_name: "dateTo")
end

defmodule AdminService.V1.ExportAuditLogsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.ExportAuditLogsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.ExportReportRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.ExportReportRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:report_type, 1, type: :string, json_name: "reportType")
  field(:date_from, 2, type: :string, json_name: "dateFrom")
  field(:date_to, 3, type: :string, json_name: "dateTo")
  field(:format, 4, type: :string)
end

defmodule AdminService.V1.ExportReportResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.ExportReportResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.GenerateSalesReportRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GenerateSalesReportRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:date_from, 1, type: :string, json_name: "dateFrom")
  field(:date_to, 2, type: :string, json_name: "dateTo")
  field(:report_type, 3, type: :string, json_name: "reportType")
end

defmodule AdminService.V1.GenerateSalesReportResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GenerateSalesReportResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.GenerateUserReportRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GenerateUserReportRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:date_from, 1, type: :string, json_name: "dateFrom")
  field(:date_to, 2, type: :string, json_name: "dateTo")
end

defmodule AdminService.V1.GenerateUserReportResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GenerateUserReportResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.GetAllShopsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetAllShopsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:approval_status, 1, type: :string, json_name: "approvalStatus")
  field(:status_filter, 2, type: :string, json_name: "statusFilter")
  field(:category, 3, type: :string)
  field(:date_from, 4, type: :string, json_name: "dateFrom")
  field(:date_to, 5, type: :string, json_name: "dateTo")
  field(:sort_by, 6, type: :string, json_name: "sortBy")
  field(:page, 7, type: :int32)
  field(:page_size, 8, type: :int32, json_name: "pageSize")
end

defmodule AdminService.V1.GetAllShopsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetAllShopsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.GetAllUsersRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetAllUsersRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:role_filter, 1, type: :string, json_name: "roleFilter")
  field(:status_filter, 2, type: :string, json_name: "statusFilter")
  field(:email_verified, 3, type: :bool, json_name: "emailVerified")
  field(:date_from, 4, type: :string, json_name: "dateFrom")
  field(:date_to, 5, type: :string, json_name: "dateTo")
  field(:sort_by, 6, type: :string, json_name: "sortBy")
  field(:page, 7, type: :int32)
  field(:page_size, 8, type: :int32, json_name: "pageSize")
end

defmodule AdminService.V1.GetAllUsersResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetAllUsersResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.GetAuditLogsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetAuditLogsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:operation_type, 1,
    type: AdminService.V1.OperationType,
    enum: true,
    json_name: "operationType"
  )

  field(:operator_id, 2, type: :string, json_name: "operatorId")
  field(:target_type, 3, type: :string, json_name: "targetType")
  field(:date_from, 4, type: :string, json_name: "dateFrom")
  field(:date_to, 5, type: :string, json_name: "dateTo")
  field(:page, 6, type: :int32)
  field(:page_size, 7, type: :int32, json_name: "pageSize")
end

defmodule AdminService.V1.GetAuditLogsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetAuditLogsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.GetCategoriesResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetCategoriesResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:categories, 2, repeated: true, type: AdminService.V1.Category)
end

defmodule AdminService.V1.GetDashboardMetricsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetDashboardMetricsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:date, 1, type: :string)
end

defmodule AdminService.V1.GetDashboardMetricsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetDashboardMetricsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.GetForbiddenWordsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetForbiddenWordsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:context, 1, type: AdminService.V1.WordContext, enum: true)
end

defmodule AdminService.V1.GetForbiddenWordsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetForbiddenWordsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.GetPendingShopsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetPendingShopsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:page, 1, type: :int32)
  field(:page_size, 2, type: :int32, json_name: "pageSize")
end

defmodule AdminService.V1.GetPendingShopsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetPendingShopsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.GetSalesChartRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetSalesChartRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:date_from, 1, type: :string, json_name: "dateFrom")
  field(:date_to, 2, type: :string, json_name: "dateTo")
  field(:group_by, 3, type: :string, json_name: "groupBy")
end

defmodule AdminService.V1.GetSalesChartResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetSalesChartResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.GetServiceHealthResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetServiceHealthResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:services, 2, repeated: true, type: AdminService.V1.ServiceHealthCheck)
end

defmodule AdminService.V1.GetSystemSettingsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetSystemSettingsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:settings, 2, repeated: true, type: AdminService.V1.SystemSettings)
end

defmodule AdminService.V1.GetUserDetailRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetUserDetailRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
end

defmodule AdminService.V1.GetUserDetailResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetUserDetailResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.RejectShopRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.RejectShopRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:admin_id, 2, type: :string, json_name: "adminId")
  field(:reason, 3, type: :string)
end

defmodule AdminService.V1.RejectShopResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.RejectShopResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.ServiceHealthCheck do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.ServiceHealthCheck",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:service_name, 2, type: :string, json_name: "serviceName")
  field(:status, 3, type: AdminService.V1.HealthStatus, enum: true)
  field(:response_time_ms, 4, type: :int32, json_name: "responseTimeMs")
  field(:error_message, 5, type: :string, json_name: "errorMessage")
  field(:checked_at, 6, type: Google.Protobuf.Timestamp, json_name: "checkedAt")
end

defmodule AdminService.V1.SuspendShopRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.SuspendShopRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:admin_id, 2, type: :string, json_name: "adminId")
  field(:reason, 3, type: :string)
end

defmodule AdminService.V1.SuspendShopResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.SuspendShopResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.SuspendUserRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.SuspendUserRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:admin_id, 2, type: :string, json_name: "adminId")
  field(:reason, 3, type: :string)
end

defmodule AdminService.V1.SuspendUserResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.SuspendUserResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.SystemSettings do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.SystemSettings",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:setting_key, 2, type: :string, json_name: "settingKey")
  field(:setting_value, 3, type: :string, json_name: "settingValue")
  field(:setting_type, 4, type: AdminService.V1.SettingType, enum: true, json_name: "settingType")
  field(:description, 5, type: :string)
  field(:updated_by, 6, type: :string, json_name: "updatedBy")
  field(:created_at, 7, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 8, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule AdminService.V1.UpdateCategoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.UpdateCategoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:admin_id, 1, type: :string, json_name: "adminId")
  field(:category_id, 2, type: :string, json_name: "categoryId")
  field(:name, 3, type: :string)
  field(:display_order, 4, type: :int32, json_name: "displayOrder")
end

defmodule AdminService.V1.UpdateCategoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.UpdateCategoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.UpdateSystemSettingRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.UpdateSystemSettingRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:admin_id, 1, type: :string, json_name: "adminId")
  field(:setting_key, 2, type: :string, json_name: "settingKey")
  field(:setting_value, 3, type: :string, json_name: "settingValue")
end

defmodule AdminService.V1.UpdateSystemSettingResponse do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.UpdateSystemSettingResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AdminService.V1.GetCategoriesRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetCategoriesRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3
end

defmodule AdminService.V1.GetServiceHealthRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetServiceHealthRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3
end

defmodule AdminService.V1.GetSystemSettingsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "admin_service.v1.GetSystemSettingsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3
end

defmodule AdminService.V1.AdminService.Service do
  @moduledoc false

  use GRPC.Service,
    name: "admin_service.v1.AdminService",
    protoc_gen_elixir_version: "0.16.0"

  rpc(:GetAllUsers, AdminService.V1.GetAllUsersRequest, AdminService.V1.GetAllUsersResponse)
  rpc(:GetUserDetail, AdminService.V1.GetUserDetailRequest, AdminService.V1.GetUserDetailResponse)

  rpc(
    :ChangeUserRole,
    AdminService.V1.ChangeUserRoleRequest,
    AdminService.V1.ChangeUserRoleResponse
  )

  rpc(:SuspendUser, AdminService.V1.SuspendUserRequest, AdminService.V1.SuspendUserResponse)
  rpc(:ActivateUser, AdminService.V1.ActivateUserRequest, AdminService.V1.ActivateUserResponse)

  rpc(
    :GetPendingShops,
    AdminService.V1.GetPendingShopsRequest,
    AdminService.V1.GetPendingShopsResponse
  )

  rpc(:ApproveShop, AdminService.V1.ApproveShopRequest, AdminService.V1.ApproveShopResponse)
  rpc(:RejectShop, AdminService.V1.RejectShopRequest, AdminService.V1.RejectShopResponse)
  rpc(:GetAllShops, AdminService.V1.GetAllShopsRequest, AdminService.V1.GetAllShopsResponse)
  rpc(:SuspendShop, AdminService.V1.SuspendShopRequest, AdminService.V1.SuspendShopResponse)
  rpc(:ActivateShop, AdminService.V1.ActivateShopRequest, AdminService.V1.ActivateShopResponse)

  rpc(
    :GetSystemSettings,
    AdminService.V1.GetSystemSettingsRequest,
    AdminService.V1.GetSystemSettingsResponse
  )

  rpc(
    :UpdateSystemSetting,
    AdminService.V1.UpdateSystemSettingRequest,
    AdminService.V1.UpdateSystemSettingResponse
  )

  rpc(:GetCategories, AdminService.V1.GetCategoriesRequest, AdminService.V1.GetCategoriesResponse)

  rpc(
    :CreateCategory,
    AdminService.V1.CreateCategoryRequest,
    AdminService.V1.CreateCategoryResponse
  )

  rpc(
    :UpdateCategory,
    AdminService.V1.UpdateCategoryRequest,
    AdminService.V1.UpdateCategoryResponse
  )

  rpc(
    :DeleteCategory,
    AdminService.V1.DeleteCategoryRequest,
    AdminService.V1.DeleteCategoryResponse
  )

  rpc(
    :GetForbiddenWords,
    AdminService.V1.GetForbiddenWordsRequest,
    AdminService.V1.GetForbiddenWordsResponse
  )

  rpc(
    :AddForbiddenWord,
    AdminService.V1.AddForbiddenWordRequest,
    AdminService.V1.AddForbiddenWordResponse
  )

  rpc(
    :DeleteForbiddenWord,
    AdminService.V1.DeleteForbiddenWordRequest,
    AdminService.V1.DeleteForbiddenWordResponse
  )

  rpc(
    :GetDashboardMetrics,
    AdminService.V1.GetDashboardMetricsRequest,
    AdminService.V1.GetDashboardMetricsResponse
  )

  rpc(:GetSalesChart, AdminService.V1.GetSalesChartRequest, AdminService.V1.GetSalesChartResponse)

  rpc(
    :GetServiceHealth,
    AdminService.V1.GetServiceHealthRequest,
    AdminService.V1.GetServiceHealthResponse
  )

  rpc(:GetAuditLogs, AdminService.V1.GetAuditLogsRequest, AdminService.V1.GetAuditLogsResponse)

  rpc(
    :ExportAuditLogs,
    AdminService.V1.ExportAuditLogsRequest,
    AdminService.V1.ExportAuditLogsResponse
  )

  rpc(
    :GenerateSalesReport,
    AdminService.V1.GenerateSalesReportRequest,
    AdminService.V1.GenerateSalesReportResponse
  )

  rpc(
    :GenerateUserReport,
    AdminService.V1.GenerateUserReportRequest,
    AdminService.V1.GenerateUserReportResponse
  )

  rpc(:ExportReport, AdminService.V1.ExportReportRequest, AdminService.V1.ExportReportResponse)
end

defmodule AdminService.V1.AdminService.Stub do
  @moduledoc false

  use GRPC.Stub, service: AdminService.V1.AdminService.Service
end
