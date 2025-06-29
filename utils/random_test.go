package utils

import (
	"testing"
	"unicode"
)

func TestGenName(t *testing.T) {
	// 测试默认长度
	t.Run("默认长度", func(t *testing.T) {
		result := GenName(8)

		if len(result) != 8 {
			t.Errorf("GenName(8)生成的字符串长度 = %d, 期望 8", len(result))
		}

		// 检查字符是否都是字母或数字
		for _, r := range result {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				t.Errorf("GenName(8)生成的字符串包含非字母数字字符: %q", r)
			}
		}
	})

	// 测试自定义长度
	testLengths := []int{1, 5, 10, 20}
	for _, length := range testLengths {
		t.Run("自定义长度", func(t *testing.T) {
			result := GenName(length)

			if len(result) != length {
				t.Errorf("GenName(%d)生成的字符串长度 = %d, 期望 %d", length, len(result), length)
			}
		})
	}

	// 测试生成多个随机字符串，确保它们是唯一的
	t.Run("唯一性", func(t *testing.T) {
		count := 100
		generated := make(map[string]bool)

		for i := 0; i < count; i++ {
			result := GenName(8)

			// 检查是否重复
			if generated[result] {
				t.Errorf("GenName(8)生成了重复的字符串: %s", result)
			}

			generated[result] = true
		}
	})

	// 测试错误处理 - 负数长度
	t.Run("负数长度", func(t *testing.T) {
		result := GenName(-1)

		// 应该返回默认长度8的字符串
		if len(result) != 8 {
			t.Errorf("GenName(-1)生成的字符串长度 = %d, 期望默认长度 8", len(result))
		}
	})
}
