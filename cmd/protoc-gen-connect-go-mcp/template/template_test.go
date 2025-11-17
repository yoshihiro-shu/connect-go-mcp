package template

import (
	"testing"
)

func TestEscapeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no escape needed",
			input:    "simple string",
			expected: "simple string",
		},
		{
			name:     "escape double quotes",
			input:    `string with "quotes"`,
			expected: `string with \"quotes\"`,
		},
		{
			name:     "escape newline",
			input:    "line1\nline2",
			expected: `line1\nline2`,
		},
		{
			name:     "escape multiple newlines",
			input:    "line1\nline2\nline3",
			expected: `line1\nline2\nline3`,
		},
		{
			name:     "escape backslash",
			input:    `path\to\file`,
			expected: `path\\to\\file`,
		},
		{
			name:     "escape backslash and quotes",
			input:    `path\to\"file\"`,
			expected: `path\\to\\\"file\\\"`,
		},
		{
			name:     "escape tab",
			input:    "column1\tcolumn2",
			expected: `column1\tcolumn2`,
		},
		{
			name:     "escape carriage return",
			input:    "line1\rline2",
			expected: `line1\rline2`,
		},
		{
			name:     "complex multiline with quotes",
			input:    "PurchaseOne - 申し込みをし、1つだけ商品を購入（=契約+請求+決済）する。\nMEMO: クローズド版リリースでは請求+決済はないため、契約の作成のみ行う",
			expected: `PurchaseOne - 申し込みをし、1つだけ商品を購入（=契約+請求+決済）する。\nMEMO: クローズド版リリースでは請求+決済はないため、契約の作成のみ行う`,
		},
		{
			name:     "multiline with multiple breaks",
			input:    "SyncVpnClient\n契約に設定されているライセンス追加イベントの数だけ、VPNサーバーに対してVPNクライアントを追加する。",
			expected: `SyncVpnClient\n契約に設定されているライセンス追加イベントの数だけ、VPNサーバーに対してVPNクライアントを追加する。`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeString(tt.input)
			if got != tt.expected {
				t.Errorf("escapeString() = %q, want %q", got, tt.expected)
			}
		})
	}
}
