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

		if findOutbound(outbounds, "AmyTelecom") != nil {
			t.Errorf("outbound-template 条目不应出现在最终输出中")
		}

		node := findOutbound(outbounds, "AmyHK")
		if node == nil {
			t.Fatalf("未找到生成的节点出站 AmyHK")
		}
		if node["domain_resolver"] != "dns_amy" {
			t.Errorf("节点未套用模板字段 domain_resolver，实际: %v", node["domain_resolver"])
		}
	})
}
