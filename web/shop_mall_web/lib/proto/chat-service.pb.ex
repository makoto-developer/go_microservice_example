# microservices/chat/proto から pb2ex.py で自動生成した最小スタブ。
defmodule ChatService.V1.MessageType do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "chat_service.v1.MessageType",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:MESSAGE_TYPE_UNSPECIFIED, 0)
  field(:TEXT, 1)
  field(:IMAGE, 2)
  field(:FILE, 3)
end

defmodule ChatService.V1.PresenceStatus do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "chat_service.v1.PresenceStatus",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:PRESENCE_STATUS_UNSPECIFIED, 0)
  field(:ONLINE, 1)
  field(:OFFLINE, 2)
  field(:AWAY, 3)
end

defmodule ChatService.V1.RoomStatus do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "chat_service.v1.RoomStatus",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:ROOM_STATUS_UNSPECIFIED, 0)
  field(:ACTIVE, 1)
  field(:RESOLVED, 2)
  field(:CLOSED, 3)
end

defmodule ChatService.V1.CreateChatRoomRequest do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.CreateChatRoomRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
  field(:shop_id, 2, type: :string, json_name: "shopId")
  field(:product_id, 3, type: :string, json_name: "productId")
end

defmodule ChatService.V1.CreateChatRoomResponse do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.CreateChatRoomResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ChatService.V1.GetArchivedMessagesRequest do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.GetArchivedMessagesRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:room_id, 1, type: :string, json_name: "roomId")
  field(:user_id, 2, type: :string, json_name: "userId")
  field(:date_from, 3, type: :string, json_name: "dateFrom")
  field(:date_to, 4, type: :string, json_name: "dateTo")
end

defmodule ChatService.V1.GetArchivedMessagesResponse do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.GetArchivedMessagesResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ChatService.V1.GetChatRoomDetailRequest do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.GetChatRoomDetailRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:room_id, 1, type: :string, json_name: "roomId")
  field(:user_id, 2, type: :string, json_name: "userId")
end

defmodule ChatService.V1.GetChatRoomDetailResponse do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.GetChatRoomDetailResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ChatService.V1.GetChatRoomsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.GetChatRoomsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:user_role, 2, type: :string, json_name: "userRole")
  field(:status_filter, 3, type: ChatService.V1.RoomStatus, enum: true, json_name: "statusFilter")
end

defmodule ChatService.V1.GetChatRoomsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.GetChatRoomsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ChatService.V1.GetMessagesRequest do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.GetMessagesRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:room_id, 1, type: :string, json_name: "roomId")
  field(:user_id, 2, type: :string, json_name: "userId")
  field(:page, 3, type: :int32)
  field(:page_size, 4, type: :int32, json_name: "pageSize")
end

defmodule ChatService.V1.GetMessagesResponse do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.GetMessagesResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ChatService.V1.GetUserPresenceRequest do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.GetUserPresenceRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
end

defmodule ChatService.V1.GetUserPresenceResponse do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.GetUserPresenceResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ChatService.V1.MarkMessagesAsReadRequest do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.MarkMessagesAsReadRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:room_id, 1, type: :string, json_name: "roomId")
  field(:user_id, 2, type: :string, json_name: "userId")
  field(:message_ids, 3, repeated: true, type: :string, json_name: "messageIds")
end

defmodule ChatService.V1.MarkMessagesAsReadResponse do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.MarkMessagesAsReadResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ChatService.V1.SearchMessagesRequest do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.SearchMessagesRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:room_id, 1, type: :string, json_name: "roomId")
  field(:user_id, 2, type: :string, json_name: "userId")
  field(:keyword, 3, type: :string)
  field(:date_from, 4, type: :string, json_name: "dateFrom")
  field(:date_to, 5, type: :string, json_name: "dateTo")
end

defmodule ChatService.V1.SearchMessagesResponse do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.SearchMessagesResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ChatService.V1.SendMessageRequest do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.SendMessageRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:room_id, 1, type: :string, json_name: "roomId")
  field(:sender_id, 2, type: :string, json_name: "senderId")
  field(:receiver_id, 3, type: :string, json_name: "receiverId")
  field(:message_type, 4, type: ChatService.V1.MessageType, enum: true, json_name: "messageType")
  field(:content, 5, type: :string)
  field(:file_url, 6, type: :string, json_name: "fileUrl")
  field(:file_name, 7, type: :string, json_name: "fileName")
  field(:file_size, 8, type: :int32, json_name: "fileSize")
end

defmodule ChatService.V1.SendMessageResponse do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.SendMessageResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ChatService.V1.UpdatePresenceRequest do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.UpdatePresenceRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:status, 2, type: ChatService.V1.PresenceStatus, enum: true)
  field(:device_info, 3, type: :string, json_name: "deviceInfo")
end

defmodule ChatService.V1.UpdatePresenceResponse do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.UpdatePresenceResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ChatService.V1.UpdateRoomStatusRequest do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.UpdateRoomStatusRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:room_id, 1, type: :string, json_name: "roomId")
  field(:user_id, 2, type: :string, json_name: "userId")
  field(:new_status, 3, type: ChatService.V1.RoomStatus, enum: true, json_name: "newStatus")
end

defmodule ChatService.V1.UpdateRoomStatusResponse do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.UpdateRoomStatusResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ChatService.V1.UploadChatFileRequest do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.UploadChatFileRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:room_id, 1, type: :string, json_name: "roomId")
  field(:user_id, 2, type: :string, json_name: "userId")
  field(:file_data, 3, type: :bytes, json_name: "fileData")
  field(:file_name, 4, type: :string, json_name: "fileName")
end

defmodule ChatService.V1.UploadChatFileResponse do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.UploadChatFileResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ChatService.V1.UploadChatImageRequest do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.UploadChatImageRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:room_id, 1, type: :string, json_name: "roomId")
  field(:user_id, 2, type: :string, json_name: "userId")
  field(:image_data, 3, type: :bytes, json_name: "imageData")
end

defmodule ChatService.V1.UploadChatImageResponse do
  @moduledoc false

  use Protobuf,
    full_name: "chat_service.v1.UploadChatImageResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ChatService.V1.ChatService.Service do
  @moduledoc false

  use GRPC.Service,
    name: "chat_service.v1.ChatService",
    protoc_gen_elixir_version: "0.16.0"

  rpc(
    :CreateChatRoom,
    ChatService.V1.CreateChatRoomRequest,
    ChatService.V1.CreateChatRoomResponse
  )

  rpc(:GetChatRooms, ChatService.V1.GetChatRoomsRequest, ChatService.V1.GetChatRoomsResponse)

  rpc(
    :GetChatRoomDetail,
    ChatService.V1.GetChatRoomDetailRequest,
    ChatService.V1.GetChatRoomDetailResponse
  )

  rpc(
    :UpdateRoomStatus,
    ChatService.V1.UpdateRoomStatusRequest,
    ChatService.V1.UpdateRoomStatusResponse
  )

  rpc(:SendMessage, ChatService.V1.SendMessageRequest, ChatService.V1.SendMessageResponse)
  rpc(:GetMessages, ChatService.V1.GetMessagesRequest, ChatService.V1.GetMessagesResponse)

  rpc(
    :MarkMessagesAsRead,
    ChatService.V1.MarkMessagesAsReadRequest,
    ChatService.V1.MarkMessagesAsReadResponse
  )

  rpc(
    :SearchMessages,
    ChatService.V1.SearchMessagesRequest,
    ChatService.V1.SearchMessagesResponse
  )

  rpc(
    :UploadChatImage,
    ChatService.V1.UploadChatImageRequest,
    ChatService.V1.UploadChatImageResponse
  )

  rpc(
    :UploadChatFile,
    ChatService.V1.UploadChatFileRequest,
    ChatService.V1.UploadChatFileResponse
  )

  rpc(
    :UpdatePresence,
    ChatService.V1.UpdatePresenceRequest,
    ChatService.V1.UpdatePresenceResponse
  )

  rpc(
    :GetUserPresence,
    ChatService.V1.GetUserPresenceRequest,
    ChatService.V1.GetUserPresenceResponse
  )

  rpc(
    :GetArchivedMessages,
    ChatService.V1.GetArchivedMessagesRequest,
    ChatService.V1.GetArchivedMessagesResponse
  )
end

defmodule ChatService.V1.ChatService.Stub do
  @moduledoc false

  use GRPC.Stub, service: ChatService.V1.ChatService.Service
end
