package clash

import (
	"fmt"
	"strconv"
	"strings"
)

// getString 安全读取 map 中的字符串字段。
// 兼容 yaml 解析得到的非字符串类型（如数字、bool），统一转为去除首尾空白的字符串。
// 字段不存在或值为空时返回 "", false。
func getString(m map[string]any, key string) (string, bool) {
	v, exist := m[key]
	if !exist || v == nil {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		s = fmt.Sprintf("%v", v)
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	return s, true
}

// getInt 安全读取 map 中的整数字段。
// 兼容 int / int64 / float64 / 数字字符串。非数字时返回 0, false。
func getInt(m map[string]any, key string) (int, bool) {
	v, exist := m[key]
	if !exist || v == nil {
		return 0, false
	}
	switch n := v.(type) {
	case int:
		return n, true
	case int32:
		return int(n), true
	case int64:
		return int(n), true
	case uint:
		return int(n), true
	case uint32:
		return int(n), true
	case uint64:
		return int(n), true
	case float32:
		return int(n), true
	case float64:
		return int(n), true
	case string:
		s := strings.TrimSpace(n)
		if s == "" {
			return 0, false
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, false
		}
		return i, true
	}
	return 0, false
}

// getBool 安全读取 map 中的布尔字段。
// 兼容 bool / 字符串 "true"|"false"|"1"|"0" / 数字 0|1。
func getBool(m map[string]any, key string) (bool, bool) {
	v, exist := m[key]
	if !exist || v == nil {
		return false, false
	}
	switch b := v.(type) {
	case bool:
		return b, true
	case string:
		s := strings.ToLower(strings.TrimSpace(b))
		switch s {
		case "true", "1", "yes":
			return true, true
		case "false", "0", "no":
			return false, true
		}
	case int:
		return b != 0, true
	case int64:
		return b != 0, true
	case float64:
		return b != 0, true
	}
	return false, false
}

// getStringSlice 读取字符串数组字段。yaml 中可能是字符串（单值）或数组。
func getStringSlice(m map[string]any, key string) ([]string, bool) {
	v, exist := m[key]
	if !exist || v == nil {
		return nil, false
	}
	switch arr := v.(type) {
	case []string:
		if len(arr) == 0 {
			return nil, false
		}
		return arr, true
	case []any:
		if len(arr) == 0 {
			return nil, false
		}
		result := make([]string, 0, len(arr))
		for _, item := range arr {
			if item == nil {
				continue
			}
			s, ok := item.(string)
			if !ok {
				s = fmt.Sprintf("%v", item)
			}
			s = strings.TrimSpace(s)
			if s != "" {
				result = append(result, s)
			}
		}
		if len(result) == 0 {
			return nil, false
		}
		return result, true
	case string:
		s := strings.TrimSpace(arr)
		if s == "" {
			return nil, false
		}
		return []string{s}, true
	}
	return nil, false
}

// getMap 读取嵌套 map 字段。
func getMap(m map[string]any, key string) (map[string]any, bool) {
	v, exist := m[key]
	if !exist || v == nil {
		return nil, false
	}
	mm, ok := v.(map[string]any)
	if !ok {
		return nil, false
	}
	return mm, true
}

// singBoxUtlsFingerprints sing-box 支持的 uTLS 指纹白名单。
// 来源: sing-box.sagernet.org/configuration/shared/tls/#utls
var singBoxUtlsFingerprints = map[string]struct{}{
	"chrome":     {},
	"firefox":    {},
	"edge":       {},
	"safari":     {},
	"360":        {},
	"qq":         {},
	"ios":        {},
	"android":    {},
	"random":     {},
	"randomized": {},
}

// normalizeUtlsFingerprint 将 fingerprint 规范化为小写并校验是否在白名单内。
// 不在白名单的 fingerprint 返回 "", false（避免输出 sing-box 不识别的值）。
func normalizeUtlsFingerprint(value string) (string, bool) {
	v := strings.ToLower(strings.TrimSpace(value))
	if v == "" {
		return "", false
	}
	if _, ok := singBoxUtlsFingerprints[v]; !ok {
		return "", false
	}
	return v, true
}
