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
			slog.Error("Failed to get current working directory", "error", err)
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
			slog.Info("providers.json not found, creating a default one.", "path", filepath.Join(configHome, defaultProvidersFilename))
			if err := createDefaultProvidersFile(filepath.Join(configHome, defaultProvidersFilename)); err != nil {
				slog.Error("Failed to create default providers.json", "error", err)
				return err
			}
			// Attempt to read again after creation
			if err := viper.ReadInConfig(); err != nil {
				slog.Error("Failed to read providers.json after creating default", "error", err)
				return err
			}
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("Unable to decode providers.json into struct", "error", err)
		return err
	}

	slog.Info("providers.json loaded successfully", "path", viper.ConfigFileUsed())
	slog.Debug("Loaded configuration", "config", cfg)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		slog.Info("providers.json changed, reloading.", "event", e.Name)
		if err := viper.Unmarshal(&cfg); err != nil {
			slog.Error("Error reloading providers.json", "error", err)
		} else {
			slog.Info("providers.json reloaded successfully.")
			slog.Debug("Reloaded configuration", "config", cfg)
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
		slog.Error("Failed to create directory for default providers.json", "directory", dir, "error", err)
		return err
	}

	v := viper.New()
	v.Set("subscribes", defaultConfig.Subscribes)
	v.Set("prefix", defaultConfig.Prefix)
	v.Set("emoji", defaultConfig.Emoji)
	v.Set("exclude_protocol", defaultConfig.ExcludeProtocol)

	if err := v.WriteConfigAs(filePath); err != nil {
		slog.Error("Failed to write default providers.json using Viper", "path", filePath, "error", err)
		return err
	}

	slog.Info("Successfully created default providers.json", "path", filePath)
	return nil
}

func GetConfig() ProvidersGlobalConfig {
	return cfg
}
