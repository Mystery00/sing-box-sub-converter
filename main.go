package main

import (
	"log/slog"
	"os"
	"sing-box-sub-converter/config"
	"sing-box-sub-converter/server"
)

func main() {
	logLevel := slog.LevelInfo
	if os.Getenv("DEBUG") != "" {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)
	err := config.LoadProvidersConfig()
	if err != nil {
		slog.Error("加载providers配置失败", "error", err)
		os.Exit(1)
		return
	}

	// Log startup
	slog.Info("正在启动sing-box-sub-converter")

	// Initialize server
	srv := server.NewServer()

	// Start server
	if err := srv.Run(); err != nil {
		slog.Error("服务器错误", "error", err)
		os.Exit(1)
	}
}
