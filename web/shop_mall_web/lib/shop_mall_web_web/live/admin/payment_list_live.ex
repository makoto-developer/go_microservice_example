defmodule ShopMallWebWeb.Admin.PaymentListLive do
  @moduledoc """
  管理者向けの決済管理画面。
  全決済の一覧・状態フィルタ・詳細表示・返金操作を提供する。
  """
  use ShopMallWebWeb, :live_view

  alias ShopMallWeb.PaymentServiceClient, as: Payments

  @status_filters [
    {"all", "すべて"},
    {"pending", "支払い待ち"},
    {"succeeded", "支払い済み"},
    {"failed", "失敗"},
    {"refunded", "返金済み"}
  ]

  @impl true
  def mount(_params, session, socket) do
    {:ok,
     socket
     |> assign(:current_user_id, session["user_id"])
     |> assign(:payments, [])
     |> assign(:total_count, 0)
     |> assign(:loading, true)
     |> assign(:error, nil)
     |> assign(:status_filter, "all")
     |> assign(:status_filters, @status_filters)
     |> assign(:detail, nil)
     |> assign(:refund_target, nil)
     |> load_payments()}
  end

  defp load_payments(socket) do
    filters =
      case socket.assigns.status_filter do
        "pending" -> [status_filter: [:PAYMENT_STATUS_PENDING]]
        "succeeded" -> [status_filter: [:PAYMENT_STATUS_SUCCEEDED]]
        "failed" -> [status_filter: [:PAYMENT_STATUS_FAILED]]
        "refunded" -> [status_filter: [:PAYMENT_STATUS_REFUNDED]]
        _ -> []
      end

    case Payments.list_payments(filters) do
      {:ok, response} ->
        socket
        |> assign(:payments, response.payments)
        |> assign(:total_count, response.total_count)
        |> assign(:loading, false)
        |> assign(:error, nil)

      {:error, reason} ->
        socket
        |> assign(:loading, false)
        |> assign(:error, "決済一覧の取得に失敗しました: #{reason}")
    end
  end

  @impl true
  def handle_event("filter_status", %{"status" => status}, socket) do
    {:noreply,
     socket
     |> assign(:status_filter, status)
     |> assign(:loading, true)
     |> load_payments()}
  end

  @impl true
  def handle_event("show_detail", %{"payment-id" => payment_id}, socket) do
    case Payments.get_payment_detail(payment_id, socket.assigns.current_user_id || "", "admin") do
      {:ok, response} ->
        {:noreply, assign(socket, :detail, response.payment)}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "詳細の取得に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("close_detail", _params, socket) do
    {:noreply, assign(socket, :detail, nil)}
  end

  @impl true
  def handle_event("open_refund", %{"payment-id" => payment_id}, socket) do
    {:noreply, assign(socket, :refund_target, payment_id)}
  end

  @impl true
  def handle_event("cancel_refund", _params, socket) do
    {:noreply, assign(socket, :refund_target, nil)}
  end

  @impl true
  def handle_event("execute_refund", %{"reason" => reason}, socket) do
    payment_id = socket.assigns.refund_target
    reason = if reason == "", do: "admin refund", else: reason

    case Payments.create_refund(payment_id, "", reason) do
      {:ok, response} ->
        {:noreply,
         socket
         |> assign(:refund_target, nil)
         |> assign(:detail, nil)
         |> put_flash(:info, "返金を実行しました(返金ID: #{response.refund_id})")
         |> load_payments()}

      {:error, reason} ->
        {:noreply,
         socket
         |> assign(:refund_target, nil)
         |> put_flash(:error, "返金に失敗しました: #{reason}")}
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-50">
      <nav class="bg-gray-900 shadow-md">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16 items-center">
            <span class="text-xl font-bold text-white">🛠 管理者コンソール</span>
            <div class="flex items-center space-x-4">
              <span class="text-gray-300 text-sm font-medium">決済管理</span>
              <.link navigate="/dashboard" class="text-gray-400 hover:text-white text-sm">
                ショップへ戻る
              </.link>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between mb-6">
          <h1 class="text-2xl font-bold text-gray-900">
            決済一覧 <span class="text-sm font-normal text-gray-500">({@total_count}件)</span>
          </h1>
          <div class="flex space-x-2">
            <button
              :for={{value, label} <- @status_filters}
              phx-click="filter_status"
              phx-value-status={value}
              class={"px-3 py-1.5 rounded-full text-sm font-medium transition-colors " <>
                if(@status_filter == value,
                  do: "bg-gray-900 text-white",
                  else: "bg-white text-gray-600 border border-gray-300 hover:bg-gray-100"
                )}
            >
              {label}
            </button>
          </div>
        </div>

        <%= if @error do %>
          <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
            {@error}
          </div>
        <% end %>

        <%= if @loading do %>
          <div class="text-center py-12">
            <div class="inline-block animate-spin rounded-full h-10 w-10 border-b-2 border-gray-900">
            </div>
          </div>
        <% else %>
          <div class="bg-white shadow rounded-lg overflow-hidden">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    決済ID
                  </th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    注文ID
                  </th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    支払い方法
                  </th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                    金額
                  </th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    状態
                  </th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    作成日時
                  </th>
                  <th class="px-4 py-3"></th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <%= if @payments == [] do %>
                  <tr>
                    <td colspan="7" class="px-4 py-8 text-center text-gray-500">
                      該当する決済はありません
                    </td>
                  </tr>
                <% end %>
                <tr :for={payment <- @payments} class="hover:bg-gray-50">
                  <td class="px-4 py-3 text-sm font-mono text-gray-600">
                    {String.slice(payment.id, 0, 8)}…
                  </td>
                  <td class="px-4 py-3 text-sm font-mono text-gray-600">
                    {String.slice(payment.order_id, 0, 8)}…
                  </td>
                  <td class="px-4 py-3 text-sm text-gray-900">
                    {Payments.method_label(payment.payment_method)}
                  </td>
                  <td class="px-4 py-3 text-sm text-right font-semibold text-gray-900">
                    {Payments.format_amount(payment.amount)}
                  </td>
                  <td class="px-4 py-3">
                    <span class={"inline-flex px-2 py-0.5 rounded-full text-xs font-medium " <>
                      Payments.status_color(payment.status)}>
                      {Payments.status_label(payment.status)}
                    </span>
                  </td>
                  <td class="px-4 py-3 text-sm text-gray-500">
                    {Payments.format_timestamp(payment.created_at)}
                  </td>
                  <td class="px-4 py-3 text-right text-sm space-x-2 whitespace-nowrap">
                    <button
                      phx-click="show_detail"
                      phx-value-payment-id={payment.id}
                      class="text-blue-600 hover:text-blue-800 font-medium"
                    >
                      詳細
                    </button>
                    <button
                      :if={payment.status == :PAYMENT_STATUS_SUCCEEDED}
                      phx-click="open_refund"
                      phx-value-payment-id={payment.id}
                      class="text-red-600 hover:text-red-800 font-medium"
                    >
                      返金
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        <% end %>
        
    <!-- 詳細モーダル -->
        <div
          :if={@detail}
          class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50"
          phx-click="close_detail"
        >
          <div
            class="bg-white rounded-lg shadow-xl max-w-lg w-full mx-4 p-6"
            phx-click-away="close_detail"
          >
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-bold text-gray-900">決済詳細</h2>
              <button phx-click="close_detail" class="text-gray-400 hover:text-gray-600">✕</button>
            </div>
            <dl class="grid grid-cols-3 gap-y-3 text-sm">
              <dt class="text-gray-500">決済ID</dt>
              <dd class="col-span-2 font-mono text-gray-900 break-all">{@detail.id}</dd>
              <dt class="text-gray-500">注文ID</dt>
              <dd class="col-span-2 font-mono text-gray-900 break-all">{@detail.order_id}</dd>
              <dt class="text-gray-500">支払い方法</dt>
              <dd class="col-span-2 text-gray-900">
                {Payments.method_label(@detail.payment_method)}
              </dd>
              <dt class="text-gray-500">金額</dt>
              <dd class="col-span-2 font-semibold text-gray-900">
                {Payments.format_amount(@detail.amount)}
                <span class="text-gray-500 font-normal uppercase">{@detail.currency}</span>
              </dd>
              <dt class="text-gray-500">状態</dt>
              <dd class="col-span-2">
                <span class={"inline-flex px-2 py-0.5 rounded-full text-xs font-medium " <>
                  Payments.status_color(@detail.status)}>
                  {Payments.status_label(@detail.status)}
                </span>
              </dd>
              <dt class="text-gray-500">作成日時</dt>
              <dd class="col-span-2 text-gray-900">
                {Payments.format_timestamp(@detail.created_at)}
              </dd>
              <dt class="text-gray-500">更新日時</dt>
              <dd class="col-span-2 text-gray-900">
                {Payments.format_timestamp(@detail.updated_at)}
              </dd>
            </dl>
          </div>
        </div>
        
    <!-- 返金確認モーダル -->
        <div
          :if={@refund_target}
          class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50"
        >
          <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
            <h2 class="text-lg font-bold text-gray-900 mb-2">返金の実行</h2>
            <p class="text-sm text-gray-600 mb-4">
              決済 <span class="font-mono">{String.slice(@refund_target, 0, 8)}…</span>
              を全額返金します。この操作は取り消せません。
            </p>
            <form phx-submit="execute_refund">
              <label class="block text-sm font-medium text-gray-700 mb-1">返金理由</label>
              <input
                type="text"
                name="reason"
                placeholder="例: お客様都合による返品"
                class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm mb-4"
              />
              <div class="flex justify-end space-x-3">
                <button
                  type="button"
                  phx-click="cancel_refund"
                  class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
                >
                  キャンセル
                </button>
                <button
                  type="submit"
                  class="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700"
                >
                  返金を実行
                </button>
              </div>
            </form>
          </div>
        </div>
      </main>
    </div>
    """
  end
end
