package utils

import (
	"testing"
)

func TestRename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"香港节点1", "🇭🇰 香港节点1"},
		{"🇭🇰 香港节点2", "🇭🇰 香港节点2"},
		{"🇺🇲 美国节点", "🇺🇸 美国节点"},
		{"Taiwan Server", "🇹🇼 Taiwan Server"},
		{"Singapore Premium", "🇸🇬 Singapore Premium"},
		{"Tokyo JP 1", "🇯🇵 Tokyo JP 1"},
		{"London UK", "🇬🇧 London UK"},
		{"中国回国节点", "🇨🇳 中国回国节点"},
		{"无法匹配的节点", "无法匹配的节点"},
		{"HK Node", "🇭🇰 HK Node"},
		{"US Node", "🇺🇸 US Node"},
		{"JP Node", "🇯🇵 JP Node"},
		{"KR Node", "🇰🇷 KR Node"},
	}

	for _, test := range tests {
		result := renameNodeTagWithEmoji(test.input)
		if result != test.expected {
			t.Errorf("RenameNodeTagWithEmoji(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}
