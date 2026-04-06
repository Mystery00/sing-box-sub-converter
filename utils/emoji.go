package utils

import (
	"regexp"
	"strings"
)

var flagEmojiPattern = regexp.MustCompile(`[\x{1F1E6}-\x{1F1FF}]{2}`)

func AddNodeEmoji(name string) string {
	if name == "" {
		return name
	}

	// 检查是否已经包含Emoji标志
	if containsEmoji(name) {
		return name
	}

	return renameNodeTagWithEmoji(name)
}

func RemoveNodeEmoji(name string) string {
	if name == "" {
		return name
	}
	if !containsEmoji(name) {
		return strings.TrimSpace(name)
	}
	r := flagEmojiPattern.ReplaceAllString(name, "")
	return strings.TrimSpace(r)
}

func containsEmoji(s string) bool {
	return flagEmojiPattern.MatchString(s)
}
