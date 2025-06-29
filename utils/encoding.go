package utils

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

func Base64Decode(str string) (string, error) {
	str = strings.TrimSpace(str)
	r, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func Base64UrlDecode(str string) (string, error) {
	str = strings.TrimSpace(str)
	r, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func JsonStr(data any) string {
	marshal, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(marshal)
}
