package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	defaultProvidersFilename = "providers.json"
	envSubConfigHome         = "SUB_CONFIG_HOME"
	envDefaultSubUrl         = "SUB_URL"
)

var cfg atomic.Pointer[ProvidersGlobalConfig]

func LoadProvidersConfig() error {
	configHome, err := filepath.Abs(os.Getenv(envSubConfigHome))
	if err != nil {
		return fmt.Errorf("获取订阅配置目录失败：%w", err)
	}
	path, err := providersConfigPath(configHome)
	if err != nil {
		return err
	}
	v := viper.New()
	v.SetConfigFile(path)
	load := func() error {
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("读取订阅配置 %s 失败：%w", path, err)
		}
		var next ProvidersGlobalConfig
		if err := v.Unmarshal(&next); err != nil {
			return fmt.Errorf("解码订阅配置 %s 失败：%w", path, err)
		}
		cfg.Store(&next)
		return nil
	}
	if err := load(); err != nil {
		return err
	}
	slog.Info("订阅配置加载成功", "path", path)
	v.OnConfigChange(func(fsnotify.Event) {
		// Viper 即使读取失败也会触发回调，重新读取以确认本次修改有效。
		if err := load(); err != nil {
			slog.Error("重新加载订阅配置失败，继续使用上次配置", "error", err)
			return
		}
		slog.Info("订阅配置重新加载成功", "path", path)
	})
	v.WatchConfig()
	return nil
}

func providersConfigPath(dir string) (string, error) {
	for _, name := range []string{defaultProvidersFilename, "providers.yaml", "providers.yml"} {
		path := filepath.Join(dir, name)
		_, err := os.Stat(path)
		if err == nil {
			return path, nil
		}
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("检查订阅配置 %s 失败：%w", path, err)
		}
	}
	path := filepath.Join(dir, defaultProvidersFilename)
	return path, createDefaultProvidersFile(path)
}

func createDefaultProvidersFile(filePath string) error {
	subURL := os.Getenv(envDefaultSubUrl)
	if subURL == "" {
		subURL = "http://127.0.0.1:1080/test.txt"
	}
	defaultConfig := ProvidersGlobalConfig{
		Subscribes: []Subscription{
			{URL: subURL, Tag: "defaultSub", Prefix: "default", UserAgent: "clash"},
		},
		Prefix:          false,
		Emoji:           true,
		ExcludeProtocol: "ssr",
		ShowSubInNodes:  false,
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		slog.Error("为默认providers.json创建目录失败", "directory", dir, "error", err)
		return err
	}

	v := viper.New()
	v.Set("subscribes", defaultConfig.Subscribes)
	v.Set("prefix", defaultConfig.Prefix)
	v.Set("emoji", defaultConfig.Emoji)
	v.Set("exclude_protocol", defaultConfig.ExcludeProtocol)
	v.Set("show_sub_in_nodes", defaultConfig.ShowSubInNodes)

	if err := v.SafeWriteConfigAs(filePath); err != nil {
		slog.Error("使用Viper写入默认providers.json失败", "path", filePath, "error", err)
		return err
	}

	slog.Info("成功创建默认providers.json", "path", filePath)
	return nil
}

func GetConfig() ProvidersGlobalConfig {
	if current := cfg.Load(); current != nil {
		return *current
	}
	return ProvidersGlobalConfig{}
}
