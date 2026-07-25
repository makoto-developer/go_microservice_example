defmodule ShopMallWebWeb.ChatLive do
  @moduledoc """
  顧客⇄店舗のチャット画面。chat サービスの全 RPC を使用する:
  ルーム作成/一覧/詳細、メッセージ送信/取得/既読/検索/アーカイブ取得、
  在席状態の更新/取得、ルームの解決済み化、画像/ファイル添付。
  """
  use ShopMallWebWeb, :live_view

  alias ChatService.V1.ChatService.Stub
  alias ChatService.V1, as: PB

  @impl true
  def mount(_params, session, socket) do
    {:ok,
     socket
     |> assign(:user_id, session["user_id"] || "guest")
     |> assign(:room_id, nil)
     |> assign(:rooms_note, nil)
     |> assign(:messages_note, nil)
     |> assign(:presence_note, nil)
     |> allow_upload(:attachment, accept: :any, max_entries: 1, max_file_size: 5_000_000)
     |> load_rooms()
     |> update_presence("online")}
  end

  defp load_rooms(socket) do
    case call(fn ch ->
           Stub.get_chat_rooms(ch, %PB.GetChatRoomsRequest{
             user_id: socket.assigns.user_id,
             user_role: "customer",
             status_filter: ""
           })
         end) do
      {:ok, resp} -> assign(socket, :rooms_note, resp.message)
      {:error, _} -> socket
    end
  end

  defp update_presence(socket, status) do
    call(fn ch ->
      Stub.update_presence(ch, %PB.UpdatePresenceRequest{
        user_id: socket.assigns.user_id,
        status: status,
        device_info: "web"
      })
    end)

    socket
  end

  @impl true
  def handle_event("create_room", %{"shop_id" => shop_id, "product_id" => product_id}, socket) do
    case call(fn ch ->
           Stub.create_chat_room(ch, %PB.CreateChatRoomRequest{
             customer_id: socket.assigns.user_id,
             shop_id: shop_id,
             product_id: product_id
           })
         end) do
      {:ok, resp} ->
        {:noreply, socket |> put_flash(:info, resp.message || "ルームを作成しました") |> load_rooms()}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "ルーム作成に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("open_room", %{"room_id" => room_id}, socket) do
    with {:ok, _detail} <-
           call(fn ch ->
             Stub.get_chat_room_detail(ch, %PB.GetChatRoomDetailRequest{
               room_id: room_id,
               user_id: socket.assigns.user_id
             })
           end),
         {:ok, msgs} <-
           call(fn ch ->
             Stub.get_messages(ch, %PB.GetMessagesRequest{
               room_id: room_id,
               user_id: socket.assigns.user_id,
               page: 1,
               page_size: 50
             })
           end) do
      # 開いたら既読にする
      call(fn ch ->
        Stub.mark_messages_as_read(ch, %PB.MarkMessagesAsReadRequest{
          room_id: room_id,
          user_id: socket.assigns.user_id
        })
      end)

      {:noreply, socket |> assign(:room_id, room_id) |> assign(:messages_note, msgs.message)}
    else
      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "ルームを開けませんでした: #{reason}")}
    end
  end

  @impl true
  def handle_event("send_message", %{"content" => content}, socket) do
    case call(fn ch ->
           Stub.send_message(ch, %PB.SendMessageRequest{
             room_id: socket.assigns.room_id || "",
             sender_id: socket.assigns.user_id,
             message_type: "text",
             content: content
           })
         end) do
      {:ok, resp} ->
        {:noreply, put_flash(socket, :info, resp.message || "送信しました")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "送信に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("search_messages", %{"keyword" => keyword}, socket) do
    case call(fn ch ->
           Stub.search_messages(ch, %PB.SearchMessagesRequest{
             room_id: socket.assigns.room_id || "",
             user_id: socket.assigns.user_id,
             keyword: keyword
           })
         end) do
      {:ok, resp} -> {:noreply, assign(socket, :messages_note, "検索: #{resp.message}")}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "検索に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("load_archive", _params, socket) do
    case call(fn ch ->
           Stub.get_archived_messages(ch, %PB.GetArchivedMessagesRequest{
             room_id: socket.assigns.room_id || "",
             user_id: socket.assigns.user_id
           })
         end) do
      {:ok, resp} -> {:noreply, assign(socket, :messages_note, "アーカイブ: #{resp.message}")}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "取得に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("check_presence", %{"target_id" => target_id}, socket) do
    case call(fn ch ->
           Stub.get_user_presence(ch, %PB.GetUserPresenceRequest{user_id: target_id})
         end) do
      {:ok, resp} -> {:noreply, assign(socket, :presence_note, resp.message)}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "在席確認に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("resolve_room", _params, socket) do
    case call(fn ch ->
           Stub.update_room_status(ch, %PB.UpdateRoomStatusRequest{
             room_id: socket.assigns.room_id || "",
             user_id: socket.assigns.user_id,
             new_status: "resolved"
           })
         end) do
      {:ok, resp} ->
        {:noreply, socket |> put_flash(:info, resp.message || "解決済みにしました") |> load_rooms()}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "更新に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("validate_upload", _params, socket), do: {:noreply, socket}

  @impl true
  def handle_event("upload_attachment", _params, socket) do
    entries =
      consume_uploaded_entries(socket, :attachment, fn %{path: path}, entry ->
        {:ok, {File.read!(path), entry.client_name}}
      end)

    case entries do
      [{data, name} | _] ->
        result =
          if String.match?(name, ~r/\.(jpg|jpeg|png|gif)$/i) do
            call(fn ch ->
              Stub.upload_chat_image(ch, %PB.UploadChatImageRequest{
                room_id: socket.assigns.room_id || "",
                user_id: socket.assigns.user_id,
                image_data: data
              })
            end)
          else
            call(fn ch ->
              Stub.upload_chat_file(ch, %PB.UploadChatFileRequest{
                room_id: socket.assigns.room_id || "",
                user_id: socket.assigns.user_id,
                file_data: data,
                file_name: name
              })
            end)
          end

        case result do
          {:ok, resp} -> {:noreply, put_flash(socket, :info, resp.message || "アップロードしました")}
          {:error, reason} -> {:noreply, put_flash(socket, :error, "アップロードに失敗しました: #{reason}")}
        end

      [] ->
        {:noreply, put_flash(socket, :error, "ファイルを選択してください")}
    end
  end

  defp call(fun) do
    host = System.get_env("CHAT_SERVICE_HOST", "localhost")
    port = System.get_env("CHAT_SERVICE_PORT", "20109")

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
        {:error, "チャットサービスに接続できません: #{inspect(reason)}"}
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-50">
      <nav class="bg-white shadow-md">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16">
            <div class="flex items-center">
              <.link
                navigate="/dashboard"
                class="text-2xl font-bold text-gray-900 hover:text-gray-700"
              >
                オンラインショップモール
              </.link>
            </div>
            <div class="flex items-center space-x-4">
              <span class="text-gray-900 px-3 py-2 text-sm font-semibold">チャット</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-6">お店とチャット</h1>

        <div class="bg-white shadow rounded-lg p-4 mb-4">
          <h2 class="text-sm font-semibold text-gray-700 mb-2">ルームを作成</h2>
          <form phx-submit="create_room" class="flex space-x-2">
            <input
              type="text"
              name="shop_id"
              required
              placeholder="店舗ID"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <input
              type="text"
              name="product_id"
              placeholder="商品ID(任意)"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <button
              type="submit"
              class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
            >
              作成
            </button>
          </form>
          <div :if={@rooms_note} class="text-xs text-gray-500 mt-2">ルーム: {@rooms_note}</div>
        </div>

        <div class="bg-white shadow rounded-lg p-4 mb-4">
          <h2 class="text-sm font-semibold text-gray-700 mb-2">ルームを開く</h2>
          <form phx-submit="open_room" class="flex space-x-2">
            <input
              type="text"
              name="room_id"
              required
              placeholder="ルームID"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <button
              type="submit"
              class="px-4 py-2 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              開く(既読になります)
            </button>
          </form>
          <div :if={@messages_note} class="text-xs text-gray-500 mt-2">{@messages_note}</div>
        </div>

        <%= if @room_id do %>
          <div class="bg-white shadow rounded-lg p-4 mb-4">
            <div class="flex items-center justify-between mb-2">
              <h2 class="text-sm font-semibold text-gray-700">ルーム {@room_id}</h2>
              <button phx-click="resolve_room" class="text-xs text-green-700 hover:text-green-900">
                ✓ 解決済みにする
              </button>
            </div>
            <form phx-submit="send_message" class="flex space-x-2 mb-3">
              <input
                type="text"
                name="content"
                required
                placeholder="メッセージを入力"
                class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
              />
              <button
                type="submit"
                class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
              >
                送信
              </button>
            </form>
            <form
              phx-submit="upload_attachment"
              phx-change="validate_upload"
              class="flex items-center space-x-2 mb-3"
            >
              <.live_file_input upload={@uploads.attachment} class="text-xs" />
              <button
                type="submit"
                class="px-3 py-1.5 text-xs font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
              >
                添付を送る
              </button>
            </form>
            <div class="flex space-x-2">
              <form phx-submit="search_messages" class="flex-1 flex space-x-2">
                <input
                  type="text"
                  name="keyword"
                  required
                  placeholder="メッセージ検索"
                  class="flex-1 border border-gray-300 rounded-md px-2 py-1.5 text-xs"
                />
                <button
                  type="submit"
                  class="px-3 py-1.5 text-xs text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50"
                >
                  検索
                </button>
              </form>
              <button
                phx-click="load_archive"
                class="px-3 py-1.5 text-xs text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50"
              >
                アーカイブ表示
              </button>
            </div>
          </div>
        <% end %>

        <div class="bg-white shadow rounded-lg p-4">
          <h2 class="text-sm font-semibold text-gray-700 mb-2">在席確認</h2>
          <form phx-submit="check_presence" class="flex space-x-2">
            <input
              type="text"
              name="target_id"
              required
              placeholder="ユーザーID"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <button
              type="submit"
              class="px-4 py-2 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              確認
            </button>
          </form>
          <div :if={@presence_note} class="text-xs text-gray-500 mt-2">{@presence_note}</div>
        </div>
      </main>
    </div>
    """
  end
end
