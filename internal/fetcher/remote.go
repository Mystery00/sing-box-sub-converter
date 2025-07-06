package fetcher

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type remote struct {
	client http.Client
}

func NewRemote() Fetcher {
	return remote{
		client: http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (remote) Check(url string) bool {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	return true
}

func (r remote) Fetch(url, userAgent string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("创建新请求失败", "url", url, "error", err)
		return "", err
	}

	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	} else {
		req.Header.Set("User-Agent", "sing-box-sub-converter")
	}

	slog.Info("正在获取订阅", "url", url, "userAgent", req.Header.Get("User-Agent"))

	resp, err := r.client.Do(req)
	if err != nil {
		slog.Error("获取订阅失败", "url", url, "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("订阅请求失败，状态码", "url", url, "status_code", resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		slog.Debug("失败请求的响应体", "body", string(bodyBytes))
		return "", fmt.Errorf("failed to fetch subscription from %s: status code %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("读取订阅响应体失败", "url", url, "error", err)
		return "", err
	}

	slog.Debug("成功获取订阅内容", "url", url, "content_length", len(body))
	return string(body), nil
}
