package utils

import (
	"regexp"
	"strings"
)

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
	emojiPattern := regexp.MustCompile(`[\x{1F1E6}-\x{1F1FF}]{2}`)
	r := emojiPattern.ReplaceAllString(name, "")
	return strings.TrimSpace(r)
}

func containsEmoji(s string) bool {
	emojiPattern := regexp.MustCompile(`[\x{1F1E6}-\x{1F1FF}]{2}`)
	return emojiPattern.MatchString(s)
}
