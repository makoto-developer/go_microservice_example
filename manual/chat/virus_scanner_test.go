package chat

import (
	"context"
	"testing"
)

func TestVirusScanner_ScanFile(t *testing.T) {
	scanner := NewVirusScanner()
	ctx := context.Background()

	result, err := scanner.ScanFile(ctx, "/tmp/test.txt")
	if err != nil {
		t.Fatalf("ScanFile failed: %v", err)
	}

	if !result.Clean {
		t.Error("Expected clean file")
	}
}

func TestVirusScanner_ScanContent(t *testing.T) {
	scanner := NewVirusScanner()
	ctx := context.Background()

	content := []byte("This is test content")
	result, err := scanner.ScanContent(ctx, content)
	if err != nil {
		t.Fatalf("ScanContent failed: %v", err)
	}

	if !result.Clean {
		t.Error("Expected clean content")
	}
}

func TestVirusScanner_IsFileAllowed(t *testing.T) {
	scanner := NewVirusScanner()

	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"JPG file", "image.jpg", true},
		{"PNG file", "image.png", true},
		{"PDF file", "document.pdf", true},
		{"EXE file", "malware.exe", false},
		{"SH file", "script.sh", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scanner.IsFileAllowed(tt.filename)
			if got != tt.want {
				t.Errorf("IsFileAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVirusScanner_CheckContentSize(t *testing.T) {
	scanner := NewVirusScanner()

	tests := []struct {
		name    string
		size    int64
		wantErr bool
	}{
		{"Small file", 1024, false},
		{"Medium file", 5 * 1024 * 1024, false},
		{"Max size", 10 * 1024 * 1024, false},
		{"Too large", 11 * 1024 * 1024, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := scanner.CheckContentSize(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckContentSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
