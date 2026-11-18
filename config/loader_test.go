package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadProvidersConfig(t *testing.T) {
	const yamlConfig = "# 订阅配置\nsubscribes:\n  - tag: yaml\n    userAgent: clash\nemoji: true # 显示图标\n"
	const jsonConfig = `{"subscribes":[{"tag":"json","userAgent":"clash"}],"emoji":true}`
	for _, tc := range []struct {
		name    string
		files   map[string]string
		wantTag string
	}{
		{"旧版 JSON", map[string]string{"providers.json": jsonConfig}, "json"},
		{"带注释的 YAML", map[string]string{"providers.yaml": yamlConfig}, "yaml"},
		{"支持 yml 扩展名", map[string]string{"providers.yml": yamlConfig}, "yaml"},
		{"优先使用 JSON", map[string]string{"providers.json": jsonConfig, "providers.yaml": yamlConfig}, "json"},
		{"创建默认配置", nil, "defaultSub"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv(envSubConfigHome, dir)
			t.Setenv(envDefaultSubUrl, "https://example.com/sub")
			for name, data := range tc.files {
				if err := os.WriteFile(filepath.Join(dir, name), []byte(data), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if err := LoadProvidersConfig(); err != nil {
				t.Fatal(err)
			}
			got := GetConfig()
			if !got.Emoji || len(got.Subscribes) != 1 || got.Subscribes[0].Tag != tc.wantTag || got.Subscribes[0].UserAgent != "clash" {
				t.Fatalf("订阅配置不符合预期：%+v", got)
			}
			if tc.files == nil && got.Subscribes[0].URL != os.Getenv(envDefaultSubUrl) {
				t.Fatal("默认配置未使用环境变量指定的订阅地址")
			}
		})
	}
}

func TestInvalidProvidersConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(envSubConfigHome, dir)
	if err := os.WriteFile(filepath.Join(dir, "providers.yaml"), []byte("emoji: ["), 0600); err != nil {
		t.Fatal(err)
	}
	if err := LoadProvidersConfig(); err == nil {
		t.Fatal("配置语法错误时应返回错误")
	}
	if _, err := os.Stat(filepath.Join(dir, defaultProvidersFilename)); !os.IsNotExist(err) {
		t.Fatal("配置语法错误时不应生成默认配置")
	}
}

func TestProvidersYAMLReload(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(envSubConfigHome, dir)
	path := filepath.Join(dir, "providers.yaml")
	if err := os.WriteFile(path, []byte("emoji: true\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := LoadProvidersConfig(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("prefix: true # 启用前缀\n"), 0600); err != nil {
		t.Fatal(err)
	}
	for deadline := time.Now().Add(3 * time.Second); time.Now().Before(deadline); {
		if got := GetConfig(); got.Prefix && !got.Emoji {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("热重载后未读取到新的 YAML 配置")
}
