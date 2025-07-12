package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"sing-box-sub-converter/config"
	"sing-box-sub-converter/converter"
	"sing-box-sub-converter/template"
	"strings"
)

func Vercel(w http.ResponseWriter, r *http.Request) {
	subURL := strings.TrimPrefix(r.URL.Path, "/vercel/")
	if subURL == "" || subURL == "/" {
		w.Write([]byte("Missing subscription URL"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 获取配置模板
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var configs map[string]any
	if err := json.Unmarshal(bodyBytes, &configs); err != nil {
		slog.Error("解析模板文件失败", "error", err)
		w.Write([]byte("Failed to load template"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	subscribes := make([]config.Subscription, 0)
	subscribes = append(subscribes, config.Subscription{
		URL:       subURL,
		Tag:       "single",
		Prefix:    "",
		UserAgent: "",
	})
	// 处理订阅
	nodes, err := converter.ProcessSubscribes(subscribes)
	if err != nil {
		slog.Error("处理订阅失败", "error", err)
		w.Write([]byte("Failed to process subscribes"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 节点信息添加到模板
	finalConfig, err := template.MergeToConfig(configs, nodes)
	if err != nil {
		slog.Error("合并配置失败", "error", err)
		w.Write([]byte("Failed to merge config"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := json.Marshal(finalConfig)
	if err != nil {
		slog.Error("序列化配置失败", "error", err)
		w.Write([]byte("Failed to marshal config"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
