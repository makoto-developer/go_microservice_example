defmodule AuthService.V1.Role do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "auth_service.v1.Role",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :ROLE_UNSPECIFIED, 0
  field :CUSTOMER, 1
  field :SHOP_OWNER, 2
  field :ADMIN, 3
end

defmodule AuthService.V1.User do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.User",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :id, 1, type: :string
  field :email, 2, type: :string
  field :password_hash, 3, type: :string, json_name: "passwordHash"
  field :role, 4, type: AuthService.V1.Role, enum: true
  field :email_verified, 5, type: :bool, json_name: "emailVerified"
  field :email_verification_token, 6, type: :string, json_name: "emailVerificationToken"

  field :email_verification_expires_at, 7,
    type: Google.Protobuf.Timestamp,
    json_name: "emailVerificationExpiresAt"

  field :password_reset_token, 8, type: :string, json_name: "passwordResetToken"

  field :password_reset_expires_at, 9,
    type: Google.Protobuf.Timestamp,
    json_name: "passwordResetExpiresAt"

  field :created_at, 10, type: Google.Protobuf.Timestamp, json_name: "createdAt"
  field :updated_at, 11, type: Google.Protobuf.Timestamp, json_name: "updatedAt"
end

defmodule AuthService.V1.RefreshToken do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.RefreshToken",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :id, 1, type: :string
  field :user_id, 2, type: :string, json_name: "userId"
  field :token, 3, type: :string
  field :expires_at, 4, type: Google.Protobuf.Timestamp, json_name: "expiresAt"
  field :revoked, 5, type: :bool
  field :created_at, 6, type: Google.Protobuf.Timestamp, json_name: "createdAt"
end

defmodule AuthService.V1.UserRegistrationRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.UserRegistrationRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :email, 1, type: :string
  field :password, 2, type: :string
  field :role, 3, type: AuthService.V1.Role, enum: true
end

defmodule AuthService.V1.EmailVerificationRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.EmailVerificationRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :token, 1, type: :string
end

defmodule AuthService.V1.UserLoginRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.UserLoginRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :email, 1, type: :string
  field :password, 2, type: :string
  field :keep_logged_in, 3, type: :bool, json_name: "keepLoggedIn"
end

defmodule AuthService.V1.UserLogoutRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.UserLogoutRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :refresh_token, 1, type: :string, json_name: "refreshToken"
end

defmodule AuthService.V1.TokenRefreshRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.TokenRefreshRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :refresh_token, 1, type: :string, json_name: "refreshToken"
end

defmodule AuthService.V1.TokenVerificationRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.TokenVerificationRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :access_token, 1, type: :string, json_name: "accessToken"
end

defmodule AuthService.V1.PasswordResetRequestRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.PasswordResetRequestRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :email, 1, type: :string
end

defmodule AuthService.V1.PasswordResetRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.PasswordResetRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :token, 1, type: :string
  field :new_password, 2, type: :string, json_name: "newPassword"
end

defmodule AuthService.V1.PasswordChangeRequest do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.PasswordChangeRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :user_id, 1, type: :string, json_name: "userId"
  field :current_password, 2, type: :string, json_name: "currentPassword"
  field :new_password, 3, type: :string, json_name: "newPassword"
end

defmodule AuthService.V1.RegisterResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.RegisterResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :user_id, 1, type: :string, json_name: "userId"
  field :access_token, 2, type: :string, json_name: "accessToken"
  field :refresh_token, 3, type: :string, json_name: "refreshToken"
  field :message, 4, type: :string
end

defmodule AuthService.V1.VerifyEmailResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.VerifyEmailResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :success, 1, type: :bool
  field :message, 2, type: :string
end

defmodule AuthService.V1.LoginResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.LoginResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :user_id, 1, type: :string, json_name: "userId"
  field :access_token, 2, type: :string, json_name: "accessToken"
  field :refresh_token, 3, type: :string, json_name: "refreshToken"
  field :role, 4, type: AuthService.V1.Role, enum: true
end

defmodule AuthService.V1.LogoutResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.LogoutResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :success, 1, type: :bool
  field :message, 2, type: :string
end

defmodule AuthService.V1.RefreshTokenResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.RefreshTokenResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :access_token, 1, type: :string, json_name: "accessToken"
  field :refresh_token, 2, type: :string, json_name: "refreshToken"
end

defmodule AuthService.V1.VerifyTokenResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.VerifyTokenResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :valid, 1, type: :bool
  field :user_id, 2, type: :string, json_name: "userId"
  field :role, 3, type: AuthService.V1.Role, enum: true
end

defmodule AuthService.V1.PasswordResetResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.PasswordResetResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :success, 1, type: :bool
  field :message, 2, type: :string
end

defmodule AuthService.V1.ResetPasswordResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.ResetPasswordResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :success, 1, type: :bool
  field :message, 2, type: :string
end

defmodule AuthService.V1.ChangePasswordResponse do
  @moduledoc false

  use Protobuf,
    full_name: "auth_service.v1.ChangePasswordResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field :success, 1, type: :bool
  field :message, 2, type: :string
end

defmodule AuthService.V1.AuthService.Service do
  @moduledoc false

  use GRPC.Service, name: "auth_service.v1.AuthService", protoc_gen_elixir_version: "0.16.0"

  rpc :Register, AuthService.V1.UserRegistrationRequest, AuthService.V1.RegisterResponse

  rpc :VerifyEmail, AuthService.V1.EmailVerificationRequest, AuthService.V1.VerifyEmailResponse

  rpc :Login, AuthService.V1.UserLoginRequest, AuthService.V1.LoginResponse

  rpc :Logout, AuthService.V1.UserLogoutRequest, AuthService.V1.LogoutResponse

  rpc :RefreshToken, AuthService.V1.TokenRefreshRequest, AuthService.V1.RefreshTokenResponse

  rpc :VerifyToken, AuthService.V1.TokenVerificationRequest, AuthService.V1.VerifyTokenResponse

  rpc :RequestPasswordReset,
      AuthService.V1.PasswordResetRequestRequest,
      AuthService.V1.PasswordResetResponse

  rpc :ResetPassword, AuthService.V1.PasswordResetRequest, AuthService.V1.ResetPasswordResponse

  rpc :ChangePassword, AuthService.V1.PasswordChangeRequest, AuthService.V1.ChangePasswordResponse
end

defmodule AuthService.V1.AuthService.Stub do
  @moduledoc false

  use GRPC.Stub, service: AuthService.V1.AuthService.Service
end
