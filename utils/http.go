package utils

import (
	"net/http"
	"time"
)

// 默认的User-Agent
const DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15"

// GetResponse 发送HTTP GET请求并返回响应
// 参数:
//   - url: 请求的URL
//   - customUserAgent: 自定义的User-Agent，如果为空则使用默认值
//
// 返回:
//   - *http.Response: 如果请求成功且状态码为200，则返回响应对象；否则返回nil
func GetResponse(url string, customUserAgent string) (*http.Response, error) {
	// 创建一个HTTP客户端，设置超时时间为5秒
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置User-Agent
	userAgent := DefaultUserAgent
	if customUserAgent != "" {
		userAgent = customUserAgent
	}

	// 添加请求头
	req.Header.Set("User-Agent", userAgent)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil
	}

	return resp, nil
}
