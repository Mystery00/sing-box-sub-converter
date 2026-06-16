# 订阅源出站模板（outbound-template）实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 允许在 config 模板的 `outbounds` 中定义 `type: "outbound-template"` 条目，按订阅 tag 给生成的出站节点深度合并一组基底属性。

**Architecture:** 在 `utils` 包新增纯函数 `DeepMerge`（递归合并两个 map，override 优先）。在 `template.MergeToConfig` 开头提取并剔除 `outbound-template` 条目建立 `tag → 模板字段` 映射；末尾生成出站节点时，对 `FromSub` 命中模板的节点以模板为基底、节点字段为覆盖做深度合并。

**Tech Stack:** Go，标准库 `testing`，`map[string]any` 表示 JSON 对象。

---

### Task 1: DeepMerge 深度合并工具函数

**Files:**
- Create: `utils/merge.go`
- Test: `utils/merge_test.go`

- [ ] **Step 1: 编写失败测试**

写入 `utils/merge_test.go`：

```go
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
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./utils/ -run TestDeepMerge -v`
Expected: 编译失败，`undefined: DeepMerge`

- [ ] **Step 3: 编写最小实现**

写入 `utils/merge.go`：

```go
package utils

// DeepMerge 以 base 为基底，用 override 覆盖，返回合并后的新 map。
// override 字段优先；同名嵌套 map 递归合并；数组及标量整体覆盖。
// 不修改入参。
func DeepMerge(base, override map[string]any) map[string]any {
	result := make(map[string]any, len(base))
	for k, v := range base {
		result[k] = v
	}
	for k, ov := range override {
		if bv, exist := result[k]; exist {
			bm, bok := bv.(map[string]any)
			om, ook := ov.(map[string]any)
			if bok && ook {
				result[k] = DeepMerge(bm, om)
				continue
			}
		}
		result[k] = ov
	}
	return result
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./utils/ -run TestDeepMerge -v`
Expected: PASS（全部子测试通过）

- [ ] **Step 5: 提交**

```bash
git add utils/merge.go utils/merge_test.go
git commit -m "feat: 新增 DeepMerge 深度合并工具函数"
```

---

### Task 2: 在 MergeToConfig 中提取出站模板并应用

**Files:**
- Modify: `template/template.go`（`MergeToConfig` 函数，当前 52-152 行）
- Test: `template/template_test.go`（新建）

- [ ] **Step 1: 编写失败测试**

写入 `template/template_test.go`。注意：shadowsocks 的 `Convert2SingBox` 会执行 `node.ProxyDetail.(ShadowsocksNode)`，`ProxyDetail` 必须是合法的 `clash.ShadowsocksNode`（字段 `Method`/`Password` 已导出，可在 `template` 包外构造）。

```go
package template

import (
	"testing"

	clashpkg "sing-box-sub-converter/converter/clash"
	"sing-box-sub-converter/converter/types"
)

func findOutbound(outbounds []any, tag string) map[string]any {
	for _, o := range outbounds {
		if m, ok := o.(map[string]any); ok {
			if m["tag"] == tag {
				return m
			}
		}
	}
	return nil
}

// ssNode 构造一个合法的 shadowsocks 节点，便于 Convert2SingBox 正常工作
func ssNode(tag, fromSub string) types.ProxyNode {
	return types.ProxyNode{
		Type:        types.ProxyNodeTypeShadowsocks,
		Tag:         tag,
		Address:     "1.2.3.4",
		Port:        "8388",
		FromSub:     fromSub,
		SubType:     "clash",
		ProxyDetail: clashpkg.ShadowsocksNode{Method: "aes-128-gcm", Password: "p"},
	}
}

func TestMergeToConfigOutboundTemplate(t *testing.T) {
	t.Run("模板字段套用且模板条目不出现在输出", func(t *testing.T) {
		config := map[string]any{
			"outbounds": []any{
				map[string]any{
					"type":            "outbound-template",
					"tag":             "AmyTelecom",
					"domain_resolver": "dns_amy",
				},
				map[string]any{
					"type":      "selector",
					"tag":       "PROXY",
					"outbounds": []any{"{all}"},
				},
			},
		}
		nodes := []types.ProxyNode{ssNode("AmyHK", "AmyTelecom")}

		result, err := MergeToConfig(config, nodes)
		if err != nil {
			t.Fatalf("MergeToConfig 返回错误: %v", err)
		}
		outbounds := result["outbounds"].([]any)

		// 模板条目不应出现
		if findOutbound(outbounds, "AmyTelecom") != nil {
			t.Errorf("outbound-template 条目不应出现在最终输出中")
		}

		// 节点应被套上 domain_resolver
		node := findOutbound(outbounds, "AmyHK")
		if node == nil {
			t.Fatalf("未找到生成的节点出站 AmyHK")
		}
		if node["domain_resolver"] != "dns_amy" {
			t.Errorf("节点未套用模板字段 domain_resolver，实际: %v", node["domain_resolver"])
		}
	})
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./template/ -run TestMergeToConfigOutboundTemplate -v`
Expected: FAIL（`domain_resolver` 为 nil，或模板条目仍出现在输出中）

- [ ] **Step 3: 修改 MergeToConfig 提取模板**

在 `template/template.go` 的 `MergeToConfig` 函数体开头（取得 `outbounds` 之后、`validNodes` 循环之前）插入提取逻辑，并把后续遍历改为跳过模板条目。

在 `outbounds, ok := config["outbounds"].([]any)` 之后插入：

```go
	// 提取出站模板：type=="outbound-template" 的条目，按 tag 建立映射，并从 outbounds 中剔除
	outboundTemplates := make(map[string]map[string]any)
	filteredOutbounds := make([]any, 0, len(outbounds))
	for _, o := range outbounds {
		outbound, ok := o.(map[string]any)
		if !ok {
			filteredOutbounds = append(filteredOutbounds, o)
			continue
		}
		if outbound["type"] == "outbound-template" {
			tag, _ := outbound["tag"].(string)
			tpl := make(map[string]any)
			for k, v := range outbound {
				if k == "type" || k == "tag" {
					continue
				}
				tpl[k] = v
			}
			if tag != "" {
				outboundTemplates[tag] = tpl
			}
			continue
		}
		filteredOutbounds = append(filteredOutbounds, outbound)
	}
	outbounds = filteredOutbounds
```

- [ ] **Step 4: 修改末尾生成出站节点循环应用模板**

将 template.go 当前 136-139 行：

```go
	for _, node := range appendNodes {
		o := converter.GetParser(node.SubType).Convert2SingBoxOutbounds(node)
		outbounds = append(outbounds, o)
	}
```

改为：

```go
	for _, node := range appendNodes {
		o := converter.GetParser(node.SubType).Convert2SingBoxOutbounds(node)
		if tpl, ok := outboundTemplates[node.FromSub]; ok {
			o = utils.DeepMerge(tpl, o)
		}
		outbounds = append(outbounds, o)
	}
```

并在文件 import 块加入 `"sing-box-sub-converter/utils"`。

- [ ] **Step 5: 运行测试确认通过**

Run: `go test ./template/ -run TestMergeToConfigOutboundTemplate -v`
Expected: PASS

- [ ] **Step 6: 运行全量测试与构建**

Run: `go build ./... && go test ./...`
Expected: 构建成功，所有测试通过

- [ ] **Step 7: 提交**

```bash
git add template/template.go template/template_test.go
git commit -m "feat: 支持订阅源出站模板，按 tag 深度合并出站节点属性"
```

---

### Task 3: 补充嵌套合并与无匹配场景的测试

**Files:**
- Modify: `template/template_test.go`

- [ ] **Step 1: 追加测试用例**

在 `TestMergeToConfigOutboundTemplate` 中追加两个子测试：

```go
	t.Run("嵌套tls深度合并且节点字段优先", func(t *testing.T) {
		config := map[string]any{
			"outbounds": []any{
				map[string]any{
					"type": "outbound-template",
					"tag":  "AmyTelecom",
					"tls": map[string]any{
						"utls": map[string]any{"enabled": true, "fingerprint": "chrome"},
					},
					"server_port": 9999, // 应被节点真实端口覆盖
				},
			},
		}
		nodes := []types.ProxyNode{ssNode("AmyHK", "AmyTelecom")}
		result, err := MergeToConfig(config, nodes)
		if err != nil {
			t.Fatalf("MergeToConfig 返回错误: %v", err)
		}
		node := findOutbound(result["outbounds"].([]any), "AmyHK")
		if node == nil {
			t.Fatalf("未找到节点 AmyHK")
		}
		tls, _ := node["tls"].(map[string]any)
		utls, _ := tls["utls"].(map[string]any)
		if utls["fingerprint"] != "chrome" {
			t.Errorf("模板嵌套 tls.utls 未保留: %v", node["tls"])
		}
		// ss 节点真实端口为 8388（uint16），应覆盖模板的 9999
		if node["server_port"] != uint16(8388) {
			t.Errorf("节点端口应覆盖模板 server_port，实际: %v", node["server_port"])
		}
	})

	t.Run("无匹配模板的节点保持原样", func(t *testing.T) {
		config := map[string]any{
			"outbounds": []any{
				map[string]any{
					"type": "outbound-template",
					"tag":  "OtherSub",
					"domain_resolver": "dns_other",
				},
			},
		}
		nodes := []types.ProxyNode{ssNode("AmyHK", "AmyTelecom")}
		result, err := MergeToConfig(config, nodes)
		if err != nil {
			t.Fatalf("MergeToConfig 返回错误: %v", err)
		}
		node := findOutbound(result["outbounds"].([]any), "AmyHK")
		if node == nil {
			t.Fatalf("未找到节点 AmyHK")
		}
		if _, exist := node["domain_resolver"]; exist {
			t.Errorf("未匹配模板的节点不应被套用字段")
		}
	})
```

- [ ] **Step 2: 运行测试确认通过**

Run: `go test ./template/ -run TestMergeToConfigOutboundTemplate -v`
Expected: PASS

- [ ] **Step 3: 提交**

```bash
git add template/template_test.go
git commit -m "test: 补充出站模板嵌套合并与无匹配场景测试"
```

---

## Self-Review

- **Spec coverage：** 模板格式（Task 2 Step 3 识别 `outbound-template`/剥离 `type`+`tag`）、深度合并语义（Task 1 + Task 2 Step 4）、剔除模板条目（Task 2 Step 3）、边界（Task 3 无匹配场景）、测试（Task 1/2/3）均有对应任务。
- **Placeholder：** 无 TBD/TODO；所有代码步骤含完整代码。
- **Type 一致性：** `DeepMerge(base, override map[string]any) map[string]any` 在 Task 1 定义、Task 2 Step 4 调用一致；`outboundTemplates`、`ssNode`、`findOutbound` 命名前后一致。
- **风险点：** ss 节点 `Convert2SingBox` 对 `ProxyDetail` 做类型断言，测试统一用 `ssNode` 构造合法 `clash.ShadowsocksNode`（字段已导出，可外部构造）。
