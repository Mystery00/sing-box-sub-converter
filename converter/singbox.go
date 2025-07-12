package converter

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// parseSingBoxContent 解析SingBox格式的内容
func parseSingBoxContent(content string) ([]map[string]interface{}, error) {
	// 解析JSON
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		// 尝试移除注释
		content = removeJSONComments(content)
		if err := json.Unmarshal([]byte(content), &data); err != nil {
			return nil, fmt.Errorf("解析SingBox JSON失败: %w", err)
		}
	}

	// 提取outbounds
	outbounds, ok := data["outbounds"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("未找到有效的outbounds字段")
	}

	// 过滤无关的outbound类型
	excludedTypes := map[string]bool{
		"selector": true,
		"urltest":  true,
		"direct":   true,
		"block":    true,
		"dns":      true,
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range outbounds {
		outbound, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// 检查类型是否应该排除
		if outboundType, ok := outbound["type"].(string); ok {
			if excludedTypes[outboundType] {
				continue
			}
		}
		result = append(result, outbound)
	}
	return result, nil
}

// removeJSONComments 移除JSON中的注释
func removeJSONComments(input string) string {
	// 使用正则表达式移除//开头的注释
	commentRegex := regexp.MustCompile(`//.*`)
	return commentRegex.ReplaceAllString(input, "")
}
