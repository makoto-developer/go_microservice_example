defmodule AuthService.V1.Role do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "auth_service.v1.Role",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:ROLE_UNSPECIFIED, 0)
  field(:CUSTOMER, 1)
  field(:SHOP_OWNER, 2)
  field(:ADMIN, 3)
end

defmodule AuthService.V1.User do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.User",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:email, 2, type: :string)
  field(:password_hash, 3, type: :string, json_name: "passwordHash")
  field(:role, 4, type: AuthService.V1.Role, enum: true)
  field(:email_verified, 5, type: :bool, json_name: "emailVerified")
  field(:email_verification_token, 6, type: :string, json_name: "emailVerificationToken")

  field(:email_verification_expires_at, 7,
    type: Google.Protobuf.Timestamp,
    json_name: "emailVerificationExpiresAt"
  )

  field(:password_reset_token, 8, type: :string, json_name: "passwordResetToken")

  field(:password_reset_expires_at, 9,
    type: Google.Protobuf.Timestamp,
    json_name: "passwordResetExpiresAt"
  )

  field(:created_at, 10, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 11, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule AuthService.V1.RefreshToken do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.RefreshToken",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:user_id, 2, type: :string, json_name: "userId")
  field(:token, 3, type: :string)
  field(:expires_at, 4, type: Google.Protobuf.Timestamp, json_name: "expiresAt")
  field(:revoked, 5, type: :bool)
  field(:created_at, 6, type: Google.Protobuf.Timestamp, json_name: "createdAt")
end

defmodule AuthService.V1.UserRegistrationRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.UserRegistrationRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:email, 1, type: :string)
  field(:password, 2, type: :string)
  field(:role, 3, type: AuthService.V1.Role, enum: true)
end

defmodule AuthService.V1.EmailVerificationRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.EmailVerificationRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:token, 1, type: :string)
end

defmodule AuthService.V1.UserLoginRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.UserLoginRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:email, 1, type: :string)
  field(:password, 2, type: :string)
  field(:keep_logged_in, 3, type: :bool, json_name: "keepLoggedIn")
end

defmodule AuthService.V1.UserLogoutRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.UserLogoutRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:refresh_token, 1, type: :string, json_name: "refreshToken")
end

defmodule AuthService.V1.TokenRefreshRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.TokenRefreshRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:refresh_token, 1, type: :string, json_name: "refreshToken")
end

defmodule AuthService.V1.TokenVerificationRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.TokenVerificationRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:access_token, 1, type: :string, json_name: "accessToken")
end

defmodule AuthService.V1.PasswordResetRequestRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.PasswordResetRequestRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:email, 1, type: :string)
end

defmodule AuthService.V1.PasswordResetRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.PasswordResetRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:token, 1, type: :string)
  field(:new_password, 2, type: :string, json_name: "newPassword")
end

defmodule AuthService.V1.PasswordChangeRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.PasswordChangeRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:current_password, 2, type: :string, json_name: "currentPassword")
  field(:new_password, 3, type: :string, json_name: "newPassword")
end

defmodule AuthService.V1.RegisterResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.RegisterResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:access_token, 2, type: :string, json_name: "accessToken")
  field(:refresh_token, 3, type: :string, json_name: "refreshToken")
  field(:message, 4, type: :string)
end

defmodule AuthService.V1.VerifyEmailResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.VerifyEmailResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AuthService.V1.LoginResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.LoginResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:access_token, 2, type: :string, json_name: "accessToken")
  field(:refresh_token, 3, type: :string, json_name: "refreshToken")
  field(:role, 4, type: AuthService.V1.Role, enum: true)
end

defmodule AuthService.V1.LogoutResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.LogoutResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AuthService.V1.RefreshTokenResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.RefreshTokenResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:access_token, 1, type: :string, json_name: "accessToken")
  field(:refresh_token, 2, type: :string, json_name: "refreshToken")
end

defmodule AuthService.V1.VerifyTokenResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.VerifyTokenResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:valid, 1, type: :bool)
  field(:user_id, 2, type: :string, json_name: "userId")
  field(:role, 3, type: AuthService.V1.Role, enum: true)
end

defmodule AuthService.V1.PasswordResetResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.PasswordResetResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AuthService.V1.ResetPasswordResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.ResetPasswordResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AuthService.V1.ChangePasswordResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.ChangePasswordResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule AuthService.V1.AuthService.Service do
  @moduledoc false

  use GRPC.Service, name: "auth_service.v1.AuthService", protoc_gen_elixir_version: "0.16.0"

  rpc(:Register, AuthService.V1.UserRegistrationRequest, AuthService.V1.RegisterResponse)

  rpc(:VerifyEmail, AuthService.V1.EmailVerificationRequest, AuthService.V1.VerifyEmailResponse)

  rpc(:Login, AuthService.V1.UserLoginRequest, AuthService.V1.LoginResponse)

  rpc(:Logout, AuthService.V1.UserLogoutRequest, AuthService.V1.LogoutResponse)

  rpc(:RefreshToken, AuthService.V1.TokenRefreshRequest, AuthService.V1.RefreshTokenResponse)

  rpc(:VerifyToken, AuthService.V1.TokenVerificationRequest, AuthService.V1.VerifyTokenResponse)

  rpc(
    :RequestPasswordReset,
    AuthService.V1.PasswordResetRequestRequest,
    AuthService.V1.PasswordResetResponse
  )

  rpc(:ResetPassword, AuthService.V1.PasswordResetRequest, AuthService.V1.ResetPasswordResponse)

  rpc(
    :ChangePassword,
    AuthService.V1.PasswordChangeRequest,
    AuthService.V1.ChangePasswordResponse
  )
end

defmodule AuthService.V1.AuthService.Stub do
  @moduledoc false

  use GRPC.Stub, service: AuthService.V1.AuthService.Service
end

# ---- CustomerAuthService(顧客専用の認証)最小スタブ ----

defmodule CustomerAuth.V1.CustomerVerifyEmailRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_auth.v1.CustomerVerifyEmailRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:token, 1, type: :string)
end

defmodule CustomerAuth.V1.CustomerVerifyEmailResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_auth.v1.CustomerVerifyEmailResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule CustomerAuth.V1.CustomerRequestPasswordResetRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_auth.v1.CustomerRequestPasswordResetRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:email, 1, type: :string)
end

defmodule CustomerAuth.V1.CustomerRequestPasswordResetResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_auth.v1.CustomerRequestPasswordResetResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule CustomerAuth.V1.CustomerResetPasswordRequest do
  @moduledoc false

  use Protobuf,
    full_name: "customer_auth.v1.CustomerResetPasswordRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:token, 1, type: :string)
  field(:new_password, 2, type: :string, json_name: "newPassword")
end

defmodule CustomerAuth.V1.CustomerResetPasswordResponse do
  @moduledoc false

  use Protobuf,
    full_name: "customer_auth.v1.CustomerResetPasswordResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule CustomerAuth.V1.CustomerAuthService.Service do
  @moduledoc false

  use GRPC.Service,
    name: "customer_auth.v1.CustomerAuthService",
    protoc_gen_elixir_version: "0.16.0"

  rpc(
    :VerifyEmail,
    CustomerAuth.V1.CustomerVerifyEmailRequest,
    CustomerAuth.V1.CustomerVerifyEmailResponse
  )

  rpc(
    :RequestPasswordReset,
    CustomerAuth.V1.CustomerRequestPasswordResetRequest,
    CustomerAuth.V1.CustomerRequestPasswordResetResponse
  )

  rpc(
    :ResetPassword,
    CustomerAuth.V1.CustomerResetPasswordRequest,
    CustomerAuth.V1.CustomerResetPasswordResponse
  )
end

defmodule CustomerAuth.V1.CustomerAuthService.Stub do
  @moduledoc false

  use GRPC.Stub, service: CustomerAuth.V1.CustomerAuthService.Service
end

# ---- OwnerAuthService(店舗オーナー専用の認証)最小スタブ ----

defmodule OwnerAuth.V1.OwnerVerifyEmailRequest do
  @moduledoc false

  use Protobuf,
    full_name: "owner_auth.v1.OwnerVerifyEmailRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:token, 1, type: :string)
end

defmodule OwnerAuth.V1.OwnerVerifyEmailResponse do
  @moduledoc false

  use Protobuf,
    full_name: "owner_auth.v1.OwnerVerifyEmailResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule OwnerAuth.V1.OwnerRequestPasswordResetRequest do
  @moduledoc false

  use Protobuf,
    full_name: "owner_auth.v1.OwnerRequestPasswordResetRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:email, 1, type: :string)
end

defmodule OwnerAuth.V1.OwnerRequestPasswordResetResponse do
  @moduledoc false

  use Protobuf,
    full_name: "owner_auth.v1.OwnerRequestPasswordResetResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule OwnerAuth.V1.OwnerResetPasswordRequest do
  @moduledoc false

  use Protobuf,
    full_name: "owner_auth.v1.OwnerResetPasswordRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:token, 1, type: :string)
  field(:new_password, 2, type: :string, json_name: "newPassword")
end

defmodule OwnerAuth.V1.OwnerResetPasswordResponse do
  @moduledoc false

  use Protobuf,
    full_name: "owner_auth.v1.OwnerResetPasswordResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule OwnerAuth.V1.OwnerAuthService.Service do
  @moduledoc false

  use GRPC.Service,
    name: "owner_auth.v1.OwnerAuthService",
    protoc_gen_elixir_version: "0.16.0"

  rpc(
    :VerifyEmail,
    OwnerAuth.V1.OwnerVerifyEmailRequest,
    OwnerAuth.V1.OwnerVerifyEmailResponse
  )

  rpc(
    :RequestPasswordReset,
    OwnerAuth.V1.OwnerRequestPasswordResetRequest,
    OwnerAuth.V1.OwnerRequestPasswordResetResponse
  )

  rpc(
    :ResetPassword,
    OwnerAuth.V1.OwnerResetPasswordRequest,
    OwnerAuth.V1.OwnerResetPasswordResponse
  )
end

defmodule OwnerAuth.V1.OwnerAuthService.Stub do
  @moduledoc false

  use GRPC.Stub, service: OwnerAuth.V1.OwnerAuthService.Service
end
