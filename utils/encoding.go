package utils

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"
)

func Base64Decode(str string) ([]byte, error) {
	// URL解码
	decoded, err := url.QueryUnescape(strings.TrimSpace(str))
	if err != nil {
		return nil, err
	}

	// 补充填充字符
	padding := len(decoded) % 4
	if padding > 0 {
		decoded += strings.Repeat("=", 4-padding)
	}

	// Base64URL解码
	return base64.URLEncoding.DecodeString(decoded)
}

func JsonStr(data any) string {
	marshal, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(marshal)
}
