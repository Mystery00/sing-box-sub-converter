package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetResponse(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求头中的User-Agent
		userAgent := r.Header.Get("User-Agent")

		// 根据User-Agent返回不同的状态码
		if userAgent == "test-agent" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Custom User-Agent Test"))
		} else if userAgent == DefaultUserAgent {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Default User-Agent Test"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	// 测试用例
	tests := []struct {
		name            string
		customUserAgent string
		expectSuccess   bool
	}{
		{"使用默认User-Agent", "", true},
		{"使用自定义User-Agent", "test-agent", true},
		{"使用无效User-Agent", "invalid-agent", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := GetResponse(server.URL, tc.customUserAgent)

			if err != nil {
				t.Fatalf("GetResponse返回错误: %v", err)
			}

			if tc.expectSuccess {
				if resp == nil {
					t.Error("预期成功但得到nil响应")
				} else {
					resp.Body.Close()
				}
			} else if resp != nil {
				resp.Body.Close()
				t.Error("预期失败但得到非nil响应")
			}
		})
	}

	// 测试URL错误
	_, err := GetResponse("http://invalid-url-that-does-not-exist", "")
	if err == nil {
		t.Error("对无效URL的请求应该返回错误")
	}
}
