package shipping

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CarrierType は配送業者タイプ
type CarrierType string

const (
	CarrierYamato     CarrierType = "yamato"
	CarrierSagawa     CarrierType = "sagawa"
	CarrierJapanPost  CarrierType = "japan_post"
)

// ShipmentStatus は配送ステータス
type ShipmentStatus string

const (
	StatusPending    ShipmentStatus = "pending"
	StatusDispatched ShipmentStatus = "dispatched"
	StatusInTransit  ShipmentStatus = "in_transit"
	StatusDelivered  ShipmentStatus = "delivered"
	StatusFailed     ShipmentStatus = "failed"
)

// TrackingInfo は追跡情報
type TrackingInfo struct {
	TrackingNumber string
	Status         ShipmentStatus
	UpdatedAt      time.Time
	Location       string
	History        []TrackingHistory
}

// TrackingHistory は追跡履歴
type TrackingHistory struct {
	Status    ShipmentStatus
	Location  string
	Timestamp time.Time
	Message   string
}

// CreateShipmentRequest は配送作成リクエスト
type CreateShipmentRequest struct {
	SenderName      string
	SenderAddress   string
	ReceiverName    string
	ReceiverAddress string
	ReceiverPhone   string
	Weight          int
	Size            string
	COD             bool
	CODAmount       int64
}

// CreateShipmentResponse は配送作成レスポンス
type CreateShipmentResponse struct {
	TrackingNumber string
	Label          []byte // 配送ラベルPDF
	CreatedAt      time.Time
}

// CarrierClient は配送業者APIのモック実装
type CarrierClient struct {
	carrier CarrierType
	apiKey  string
}

// NewCarrierClient はCarrierClientを初期化
func NewCarrierClient(carrier CarrierType, apiKey string) *CarrierClient {
	return &CarrierClient{
		carrier: carrier,
		apiKey:  apiKey,
	}
}

// CreateShipment は配送を作成（モック）
func (c *CarrierClient) CreateShipment(ctx context.Context, req CreateShipmentRequest) (*CreateShipmentResponse, error) {
	// モック実装: ダミーの追跡番号を生成
	var trackingNumber string
	switch c.carrier {
	case CarrierYamato:
		trackingNumber = fmt.Sprintf("YM%d", time.Now().Unix())
	case CarrierSagawa:
		trackingNumber = fmt.Sprintf("SG%d", time.Now().Unix())
	case CarrierJapanPost:
		trackingNumber = fmt.Sprintf("JP%d", time.Now().Unix())
	default:
		trackingNumber = fmt.Sprintf("XX%d", time.Now().Unix())
	}

	resp := &CreateShipmentResponse{
		TrackingNumber: trackingNumber,
		Label:          []byte("PDF_MOCK_DATA"), // 実際はPDFバイナリ
		CreatedAt:      time.Now(),
	}

	return resp, nil
}

// GetTracking は追跡情報を取得（モック）
func (c *CarrierClient) GetTracking(ctx context.Context, trackingNumber string) (*TrackingInfo, error) {
	// モック実装: ダミーの追跡情報を返す
	info := &TrackingInfo{
		TrackingNumber: trackingNumber,
		Status:         StatusInTransit,
		UpdatedAt:      time.Now(),
		Location:       "東京都千代田区配送センター",
		History: []TrackingHistory{
			{
				Status:    StatusDispatched,
				Location:  "大阪府大阪市集荷センター",
				Timestamp: time.Now().Add(-24 * time.Hour),
				Message:   "荷物を発送しました",
			},
			{
				Status:    StatusInTransit,
				Location:  "東京都千代田区配送センター",
				Timestamp: time.Now().Add(-6 * time.Hour),
				Message:   "配送センターに到着しました",
			},
		},
	}

	return info, nil
}

// CancelShipment は配送をキャンセル（モック）
func (c *CarrierClient) CancelShipment(ctx context.Context, trackingNumber string) error {
	// モック実装: キャンセル成功を返す
	return nil
}

// CalculateShippingFee は送料を計算（モック）
func (c *CarrierClient) CalculateShippingFee(ctx context.Context, weight int, size string, fromPostalCode string, toPostalCode string) (int64, error) {
	// モック実装: サイズと重量から簡易計算
	baseFee := int64(800)

	// サイズによる加算
	switch size {
	case "small":
		baseFee = 800
	case "medium":
		baseFee = 1200
	case "large":
		baseFee = 1800
	default:
		baseFee = 1000
	}

	// 重量による加算（1kg超過ごとに200円）
	if weight > 1000 {
		extraKg := (weight - 1000) / 1000
		baseFee += int64(extraKg) * 200
	}

	return baseFee, nil
}
