defmodule ShopMallWebWeb.Admin.ReviewModerationLive do
  @moduledoc """
  管理者のレビューモデレーション画面。
  承認待ちレビューの取得(GetPendingReviews)、承認/却下/削除
  (ApproveReview / RejectReview / DeleteReviewByAdmin)、
  通報の一覧と対応(GetReviewReports / ResolveReviewReport)を行う。
  """
  use ShopMallWebWeb, :live_view

  alias ShopMallWeb.ReviewServiceClient, as: Reviews

  @impl true
  def mount(_params, session, socket) do
    {:ok,
     socket
     |> assign(:admin_id, session["user_id"] || "admin")
     |> assign(:pending_note, nil)
     |> assign(:reports_note, nil)}
  end

  @impl true
  def handle_event("load_pending", _params, socket) do
    case Reviews.get_pending_reviews() do
      {:ok, resp} -> {:noreply, assign(socket, :pending_note, resp.message)}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "取得に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("load_reports", _params, socket) do
    case Reviews.get_review_reports(socket.assigns.admin_id) do
      {:ok, resp} -> {:noreply, assign(socket, :reports_note, resp.message)}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "取得に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("approve", %{"review_id" => review_id}, socket) do
    run(socket, Reviews.approve_review(review_id, socket.assigns.admin_id))
  end

  @impl true
  def handle_event("reject", %{"review_id" => review_id, "reason" => reason}, socket) do
    run(socket, Reviews.reject_review(review_id, socket.assigns.admin_id, reason))
  end

  @impl true
  def handle_event("delete", %{"review_id" => review_id, "reason" => reason}, socket) do
    run(socket, Reviews.delete_review_by_admin(review_id, socket.assigns.admin_id, reason))
  end

  @impl true
  def handle_event("resolve_report", %{"report_id" => report_id, "action" => action}, socket) do
    run(socket, Reviews.resolve_review_report(report_id, socket.assigns.admin_id, action, ""))
  end

  defp run(socket, result) do
    case result do
      {:ok, resp} -> {:noreply, put_flash(socket, :info, resp.message || "実行しました")}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "操作に失敗しました: #{reason}")}
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
              <.link navigate="/admin/payments" class="text-gray-400 hover:text-white text-sm">
                決済管理
              </.link>
              <.link navigate="/admin/orders" class="text-gray-400 hover:text-white text-sm">
                注文分析
              </.link>
              <span class="text-gray-300 text-sm font-medium">レビュー管理</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-6">レビューモデレーション</h1>

        <div class="grid grid-cols-2 gap-4 mb-6">
          <div class="bg-white shadow rounded-lg p-4">
            <div class="flex items-center justify-between mb-2">
              <h2 class="text-sm font-semibold text-gray-700">承認待ちレビュー</h2>
              <button phx-click="load_pending" class="text-xs text-blue-600 hover:text-blue-800">
                取得
              </button>
            </div>
            <div :if={@pending_note} class="text-xs text-gray-600">{@pending_note}</div>
          </div>
          <div class="bg-white shadow rounded-lg p-4">
            <div class="flex items-center justify-between mb-2">
              <h2 class="text-sm font-semibold text-gray-700">通報一覧</h2>
              <button phx-click="load_reports" class="text-xs text-blue-600 hover:text-blue-800">
                取得
              </button>
            </div>
            <div :if={@reports_note} class="text-xs text-gray-600">{@reports_note}</div>
          </div>
        </div>

        <div class="bg-white shadow rounded-lg p-6 mb-6">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">レビューの承認 / 却下 / 削除</h2>
          <form phx-submit="approve" class="flex items-end space-x-2 mb-3">
            <input
              type="text"
              name="review_id"
              required
              placeholder="レビューID"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <button
              type="submit"
              class="px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700"
            >
              承認
            </button>
          </form>
          <form phx-submit="reject" class="flex items-end space-x-2 mb-3">
            <input
              type="text"
              name="review_id"
              required
              placeholder="レビューID"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <input
              type="text"
              name="reason"
              required
              placeholder="却下理由"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <button
              type="submit"
              class="px-4 py-2 text-sm font-medium text-yellow-700 border border-yellow-300 rounded-md hover:bg-yellow-50"
            >
              却下
            </button>
          </form>
          <form phx-submit="delete" class="flex items-end space-x-2">
            <input
              type="text"
              name="review_id"
              required
              placeholder="レビューID"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <input
              type="text"
              name="reason"
              required
              placeholder="削除理由"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <button
              type="submit"
              class="px-4 py-2 text-sm font-medium text-red-600 border border-red-300 rounded-md hover:bg-red-50"
            >
              削除
            </button>
          </form>
        </div>

        <div class="bg-white shadow rounded-lg p-6">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">通報への対応</h2>
          <form phx-submit="resolve_report" class="flex items-end space-x-2">
            <input
              type="text"
              name="report_id"
              required
              placeholder="通報ID"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <select name="action" class="border border-gray-300 rounded-md px-2 py-2 text-sm">
              <option value="dismiss">問題なし(却下)</option>
              <option value="remove_review">レビューを削除</option>
              <option value="warn_user">投稿者に警告</option>
            </select>
            <button
              type="submit"
              class="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
            >
              対応を記録
            </button>
          </form>
        </div>
      </main>
    </div>
    """
  end
end
