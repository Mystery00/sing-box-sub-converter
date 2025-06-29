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

// LoadProvidersConfig loads the providers.json configuration.
// It also sets up a watcher to reload the config on change.
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

	// Create default config if it doesn't exist
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

	// Use Viper to write the default config to ensure consistency if we use Viper's write methods elsewhere.
	// However, Viper doesn't have a direct "write this struct as config" function.
	// So, we'll marshal it to JSON and write manually.
	// For simplicity and direct control, we'll marshal and write.

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		slog.Error("Failed to create directory for default providers.json", "directory", dir, "error", err)
		return err
	}

	// Create a temporary viper instance to write the default config
	// This is a bit roundabout, as viper is more for reading and managing existing files.
	// For creating a default file, marshalling to JSON and writing directly is often simpler.
	// However, to stick to viper for all config file interactions:
	v := viper.New()
	v.Set("subscribes", defaultConfig.Subscribes)
	v.Set("prefix", defaultConfig.Prefix)
	v.Set("emoji", defaultConfig.Emoji)
	v.Set("exclude_protocol", defaultConfig.ExcludeProtocol)

	// Attempt to write the configuration.
	// WriteConfigAs will create the file if it doesn't exist.
	if err := v.WriteConfigAs(filePath); err != nil {
		slog.Error("Failed to write default providers.json using Viper", "path", filePath, "error", err)
		return err
	}

	slog.Info("Successfully created default providers.json", "path", filePath)
	return nil
}

// GetConfig returns the current loaded configuration.
// It's good practice to return a copy or ensure cfg is accessed concurrently safely if needed,
// but for now, direct access is fine as Viper's OnConfigChange handles updates.
func GetConfig() ProvidersGlobalConfig {
	return cfg
}
