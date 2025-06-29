package utils

import (
	"testing"
)

func TestBase64Decode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "普通Base64URL编码",
			input:    "SGVsbG8gV29ybGQ",
			expected: "Hello World",
			wantErr:  false,
		},
		{
			name:     "带空格的Base64URL编码",
			input:    " SGVsbG8gV29ybGQ ",
			expected: "Hello World",
			wantErr:  false,
		},
		{
			name:     "URL编码的Base64字符串",
			input:    "SGVsbG8lMjBXb3JsZA",
			expected: "Hello World",
			wantErr:  false,
		},
		{
			name:     "无效的Base64字符串",
			input:    "!@#$%^",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Base64Decode(tc.input)

			// 检查错误
			if (err != nil) != tc.wantErr {
				t.Errorf("Base64Decode() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// 如果预期不出错，检查结果
			if !tc.wantErr && result != tc.expected {
				t.Errorf("Base64Decode() = %v, want %v", result, tc.expected)
			}
		})
	}
}
