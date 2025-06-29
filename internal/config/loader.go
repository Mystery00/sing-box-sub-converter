package config

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	defaultProvidersFilename = "providers.json"
	envSubConfigHome         = "SUB_CONFIG_HOME"
)

var cfg ProvidersGlobalConfig

func LoadProvidersConfig() error {
	configHome := os.Getenv(envSubConfigHome)
	if configHome == "" {
		cwd, err := os.Getwd()
		if err != nil {
			slog.Error("获取当前工作目录失败", "error", err)
			return err
		}
		configHome = cwd
	}

	viper.SetConfigName(strings.TrimSuffix(defaultProvidersFilename, filepath.Ext(defaultProvidersFilename)))
	viper.SetConfigType("json")
	viper.AddConfigPath(configHome)

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			slog.Info("未找到providers.json，正在创建默认文件", "path", filepath.Join(configHome, defaultProvidersFilename))
			if err := createDefaultProvidersFile(filepath.Join(configHome, defaultProvidersFilename)); err != nil {
				slog.Error("创建默认providers.json文件失败", "error", err)
				return err
			}
			// Attempt to read again after creation
			if err := viper.ReadInConfig(); err != nil {
				slog.Error("创建默认文件后读取providers.json失败", "error", err)
				return err
			}
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("无法将providers.json解码为结构体", "error", err)
		return err
	}

	slog.Info("providers.json加载成功", "path", viper.ConfigFileUsed())
	slog.Debug("已加载配置", "config", cfg)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		slog.Info("providers.json已更改，正在重新加载", "event", e.Name)
		if err := viper.Unmarshal(&cfg); err != nil {
			slog.Error("重新加载providers.json出错", "error", err)
		} else {
			slog.Info("providers.json重新加载成功")
			slog.Debug("已重新加载配置", "config", cfg)
		}
	})

	return nil
}

func createDefaultProvidersFile(filePath string) error {
	defaultConfig := ProvidersGlobalConfig{
		Subscribes: []Subscription{
			{URL: "http://127.0.0.1:1080/test1.txt", Tag: "test1", Prefix: "test1", UserAgent: "clash"},
			{URL: "http://127.0.0.1:1080/test2.txt", Tag: "test2", Prefix: "test2", UserAgent: ""},
		},
		Prefix:          true,
		Emoji:           true,
		ExcludeProtocol: "ssr",
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

	if err := v.WriteConfigAs(filePath); err != nil {
		slog.Error("使用Viper写入默认providers.json失败", "path", filePath, "error", err)
		return err
	}

	slog.Info("成功创建默认providers.json", "path", filePath)
	return nil
}

func GetConfig() ProvidersGlobalConfig {
	return cfg
}
