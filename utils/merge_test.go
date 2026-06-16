package utils

import (
	"reflect"
	"testing"
)

func TestDeepMerge(t *testing.T) {
	t.Run("标量覆盖", func(t *testing.T) {
		base := map[string]any{"a": 1, "b": 2}
		override := map[string]any{"b": 3, "c": 4}
		got := DeepMerge(base, override)
		want := map[string]any{"a": 1, "b": 3, "c": 4}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("DeepMerge = %v, 期望 %v", got, want)
		}
	})

	t.Run("嵌套map递归合并", func(t *testing.T) {
		base := map[string]any{
			"tls": map[string]any{"utls": map[string]any{"enabled": true, "fingerprint": "chrome"}},
		}
		override := map[string]any{
			"tls": map[string]any{"enabled": true, "server_name": "a.com"},
		}
		got := DeepMerge(base, override)
		want := map[string]any{
			"tls": map[string]any{
				"enabled":     true,
				"server_name": "a.com",
				"utls":        map[string]any{"enabled": true, "fingerprint": "chrome"},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("DeepMerge = %v, 期望 %v", got, want)
		}
	})

	t.Run("数组整体覆盖", func(t *testing.T) {
		base := map[string]any{"alpn": []any{"h2", "http/1.1"}}
		override := map[string]any{"alpn": []any{"h3"}}
		got := DeepMerge(base, override)
		want := map[string]any{"alpn": []any{"h3"}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("DeepMerge = %v, 期望 %v", got, want)
		}
	})

	t.Run("空base", func(t *testing.T) {
		got := DeepMerge(map[string]any{}, map[string]any{"a": 1})
		want := map[string]any{"a": 1}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("DeepMerge = %v, 期望 %v", got, want)
		}
	})

	t.Run("空override", func(t *testing.T) {
		got := DeepMerge(map[string]any{"a": 1}, map[string]any{})
		want := map[string]any{"a": 1}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("DeepMerge = %v, 期望 %v", got, want)
		}
	})

	t.Run("不修改入参", func(t *testing.T) {
		base := map[string]any{"a": 1}
		override := map[string]any{"a": 2}
		DeepMerge(base, override)
		if base["a"] != 1 {
			t.Errorf("base 被修改了: %v", base)
		}
	})
}
