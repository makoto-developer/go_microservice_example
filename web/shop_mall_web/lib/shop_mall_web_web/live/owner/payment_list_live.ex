defmodule ShopMallWebWeb.Owner.PaymentListLive do
  @moduledoc """
  加盟店(店舗オーナー)向けの決済確認画面。
  自店舗宛の決済一覧の確認と、代引き注文の集金確定(配達完了報告)を行う。

  注: このサンプルでは payments テーブルに shop_id が無いため、
  店舗によるサーバー側の絞り込みは行っていない(全件が見える)。
  """
  use ShopMallWebWeb, :live_view

  alias ShopMallWeb.PaymentServiceClient, as: Payments

  @method_filters [
    {"all", "すべて"},
    {"cod", "代金引換"},
    {"card", "クレジットカード"}
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
     |> assign(:method_filter, "all")
     |> assign(:method_filters, @method_filters)
     |> load_payments()}
  end

  defp load_payments(socket) do
    filters =
      case socket.assigns.method_filter do
        "cod" -> [payment_method: :CASH_ON_DELIVERY]
        "card" -> [payment_method: :CREDIT_CARD]
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
  def handle_event("filter_method", %{"method" => method}, socket) do
    {:noreply,
     socket
     |> assign(:method_filter, method)
     |> assign(:loading, true)
     |> load_payments()}
  end

  @impl true
  def handle_event(
        "confirm_cod",
        %{"payment-id" => payment_id, "order-id" => order_id},
        socket
      ) do
    case Payments.confirm_cod_payment(payment_id, order_id) do
      {:ok, response} ->
        {:noreply,
         socket
         |> put_flash(:info, "集金を確定しました: #{response.message}")
         |> load_payments()}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "集金の確定に失敗しました: #{reason}")}
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-50">
      <nav class="bg-emerald-800 shadow-md">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16 items-center">
            <span class="text-xl font-bold text-white">🏪 加盟店ポータル</span>
            <div class="flex items-center space-x-4">
              <.link
                navigate="/owner/dashboard"
                class="text-emerald-200 hover:text-white text-sm font-medium"
              >
                ダッシュボード
              </.link>
              <.link
                navigate="/owner/products"
                class="text-emerald-200 hover:text-white text-sm font-medium"
              >
                商品管理
              </.link>
              <span class="text-white text-sm font-medium">決済確認</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between mb-2">
          <h1 class="text-2xl font-bold text-gray-900">
            決済確認 <span class="text-sm font-normal text-gray-500">({@total_count}件)</span>
          </h1>
          <div class="flex space-x-2">
            <button
              :for={{value, label} <- @method_filters}
              phx-click="filter_method"
              phx-value-method={value}
              class={"px-3 py-1.5 rounded-full text-sm font-medium transition-colors " <>
                if(@method_filter == value,
                  do: "bg-emerald-700 text-white",
                  else: "bg-white text-gray-600 border border-gray-300 hover:bg-gray-100"
                )}
            >
              {label}
            </button>
          </div>
        </div>
        <p class="text-sm text-gray-500 mb-6">
          代金引換の注文は、商品のお届け完了後に「集金確定」を押してください。
        </p>

        <%= if @error do %>
          <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
            {@error}
          </div>
        <% end %>

        <%= if @loading do %>
          <div class="text-center py-12">
            <div class="inline-block animate-spin rounded-full h-10 w-10 border-b-2 border-emerald-700">
            </div>
          </div>
        <% else %>
          <div class="bg-white shadow rounded-lg overflow-hidden">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
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
                    受注日時
                  </th>
                  <th class="px-4 py-3"></th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <%= if @payments == [] do %>
                  <tr>
                    <td colspan="6" class="px-4 py-8 text-center text-gray-500">
                      該当する決済はありません
                    </td>
                  </tr>
                <% end %>
                <tr :for={payment <- @payments} class="hover:bg-gray-50">
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
                  <td class="px-4 py-3 text-right whitespace-nowrap">
                    <button
                      :if={
                        payment.payment_method == :CASH_ON_DELIVERY and
                          payment.status == :PAYMENT_STATUS_PENDING
                      }
                      phx-click="confirm_cod"
                      phx-value-payment-id={payment.id}
                      phx-value-order-id={payment.order_id}
                      class="px-3 py-1.5 text-sm font-medium text-white bg-emerald-600 rounded-md hover:bg-emerald-700"
                    >
                      集金確定
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        <% end %>
      </main>
    </div>
    """
  end
end
