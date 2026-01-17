package chat

import (
	"context"
	"fmt"
)

// VirusScanner はウイルススキャナ（モック）
type VirusScanner struct{}

// NewVirusScanner はVirusScannerを初期化
func NewVirusScanner() *VirusScanner {
	return &VirusScanner{}
}

// ScanResult はスキャン結果
type ScanResult struct {
	Clean     bool
	VirusName string
	ThreatLevel string
}

// ScanFile はファイルをスキャン（モック）
func (s *VirusScanner) ScanFile(ctx context.Context, filePath string) (*ScanResult, error) {
	// モック実装: 常にクリーンとする
	// 実際はClamAVなどのウイルススキャナーを使用
	fmt.Printf("[VIRUS SCAN MOCK] Scanning file: %s\n", filePath)

	return &ScanResult{
		Clean:       true,
		VirusName:   "",
		ThreatLevel: "",
	}, nil
}

// ScanContent はコンテンツをスキャン（モック）
func (s *VirusScanner) ScanContent(ctx context.Context, content []byte) (*ScanResult, error) {
	// モック実装: 常にクリーンとする
	fmt.Printf("[VIRUS SCAN MOCK] Scanning content: %d bytes\n", len(content))

	return &ScanResult{
		Clean:       true,
		VirusName:   "",
		ThreatLevel: "",
	}, nil
}

// IsFileAllowed はファイル拡張子をチェック
func (s *VirusScanner) IsFileAllowed(filename string) bool {
	// 許可する拡張子のホワイトリスト
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".pdf":  true,
		".txt":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".zip":  true,
	}

	// 拡張子を抽出
	ext := ""
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			ext = filename[i:]
			break
		}
	}

	return allowedExtensions[ext]
}

// CheckContentSize はコンテンツサイズをチェック
func (s *VirusScanner) CheckContentSize(size int64) error {
	const maxFileSize = 10 * 1024 * 1024 // 10MB

	if size > maxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size: %d bytes", maxFileSize)
	}

	return nil
}
