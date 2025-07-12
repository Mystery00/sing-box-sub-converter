package fetcher

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
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

func (r remote) Fetch(url, userAgent string) (string, *SubInfo, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("创建新请求失败", "url", url, "error", err)
		return "", nil, err
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
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("订阅请求失败，状态码", "url", url, "status_code", resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		slog.Debug("失败请求的响应体", "body", string(bodyBytes))
		return "", nil, fmt.Errorf("failed to fetch subscription from %s: status code %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("读取订阅响应体失败", "url", url, "error", err)
		return "", nil, err
	}

	slog.Debug("成功获取订阅内容", "url", url, "content_length", len(body))
	subUserInfoStr := resp.Header.Get("subscription-userinfo")
	var subInfo *SubInfo = nil
	if subUserInfoStr != "" {
		subInfo = &SubInfo{}
		// 解析subscription-userinfo头
		parts := strings.Split(subUserInfoStr, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			kv := strings.Split(part, "=")
			if len(kv) != 2 {
				continue
			}
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			// 将字符串值转换为int64并设置到subInfo结构体中
			atoi, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			switch key {
			case "upload":
				subInfo.Upload = int64(atoi)
				break
			case "download":
				subInfo.Download = int64(atoi)
				break
			case "total":
				subInfo.Total = int64(atoi)
				break
			case "expire":
				subInfo.Expire = int64(atoi)
				break
			}
		}

	}
	return string(body), subInfo, nil
}
