package shipping

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidPostalCode = errors.New("invalid postal code format")
	ErrInvalidAddress    = errors.New("invalid address")
)

// AddressNormalizer は住所正規化ロジック
type AddressNormalizer struct{}

// NewAddressNormalizer はAddressNormalizerを初期化
func NewAddressNormalizer() *AddressNormalizer {
	return &AddressNormalizer{}
}

// NormalizedAddress は正規化された住所
type NormalizedAddress struct {
	PostalCode string
	Prefecture string
	City       string
	Address1   string
	Address2   string
}

// NormalizePostalCode は郵便番号を正規化
func (n *AddressNormalizer) NormalizePostalCode(postalCode string) (string, error) {
	// ハイフンを削除
	normalized := strings.ReplaceAll(postalCode, "-", "")
	normalized = strings.ReplaceAll(normalized, "−", "")
	normalized = strings.ReplaceAll(normalized, "ー", "")

	// 数字7桁のパターンをチェック
	if matched, _ := regexp.MatchString(`^\d{7}$`, normalized); !matched {
		return "", ErrInvalidPostalCode
	}

	// ハイフン付きフォーマットに変換
	return normalized[:3] + "-" + normalized[3:], nil
}

// Normalize は住所全体を正規化（モック）
func (n *AddressNormalizer) Normalize(postalCode string, prefecture string, city string, address1 string, address2 string) (*NormalizedAddress, error) {
	// 郵便番号の正規化
	normalizedPostalCode, err := n.NormalizePostalCode(postalCode)
	if err != nil {
		return nil, err
	}

	// 都道府県名の正規化
	normalizedPrefecture := n.normalizePrefecture(prefecture)

	// 市区町村名の正規化
	normalizedCity := strings.TrimSpace(city)

	// 番地の正規化
	normalizedAddress1 := strings.TrimSpace(address1)
	normalizedAddress2 := strings.TrimSpace(address2)

	return &NormalizedAddress{
		PostalCode: normalizedPostalCode,
		Prefecture: normalizedPrefecture,
		City:       normalizedCity,
		Address1:   normalizedAddress1,
		Address2:   normalizedAddress2,
	}, nil
}

// normalizePrefecture は都道府県名を正規化
func (n *AddressNormalizer) normalizePrefecture(prefecture string) string {
	// 「都」「道」「府」「県」を統一
	normalized := strings.TrimSpace(prefecture)

	// 既に都道府県がついている場合はそのまま
	if strings.HasSuffix(normalized, "都") ||
		strings.HasSuffix(normalized, "道") ||
		strings.HasSuffix(normalized, "府") ||
		strings.HasSuffix(normalized, "県") {
		return normalized
	}

	// 都道府県を補完（簡易実装）
	prefectureMap := map[string]string{
		"東京":   "東京都",
		"北海道":  "北海道",
		"大阪":   "大阪府",
		"京都":   "京都府",
		"神奈川":  "神奈川県",
		"埼玉":   "埼玉県",
		"千葉":   "千葉県",
		// ... 他の都道府県も同様
	}

	if fullName, exists := prefectureMap[normalized]; exists {
		return fullName
	}

	// 見つからない場合は「県」を付与
	return normalized + "県"
}

// ValidateAddress は住所を検証
func (n *AddressNormalizer) ValidateAddress(address *NormalizedAddress) error {
	if address.PostalCode == "" {
		return ErrInvalidAddress
	}
	if address.Prefecture == "" {
		return ErrInvalidAddress
	}
	if address.City == "" {
		return ErrInvalidAddress
	}
	if address.Address1 == "" {
		return ErrInvalidAddress
	}

	return nil
}
