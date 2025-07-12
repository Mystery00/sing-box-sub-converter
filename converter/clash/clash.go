package clash

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"sing-box-sub-converter/converter/types"
)

type ProtocolConverter interface {
	NodeType() types.ProxyNodeType
	Handle(map[string]any) bool
	Parse(map[string]any) ([]types.ProxyNode, error)
	Convert2SingBox(types.ProxyNode) map[string]any
}

type Clash struct {
	converters   []ProtocolConverter
	converterMap map[types.ProxyNodeType]ProtocolConverter
}

func NewClash() *Clash {
	converters := []ProtocolConverter{
		shadowsocks{},
		trojan{},
		hysteria2{},
		vless{},
		vmess{},
	}

	converterMap := make(map[types.ProxyNodeType]ProtocolConverter)
	for _, protocolConverter := range converters {
		converterMap[protocolConverter.NodeType()] = protocolConverter
	}
	return &Clash{
		converters,
		converterMap,
	}
}

func (Clash) SubType() string {
	return "clash"
}

func (p Clash) Parse(content, sub string) ([]types.ProxyNode, error) {
	var data map[string]any
	if err := yaml.Unmarshal([]byte(content), &data); err != nil {
		return nil, fmt.Errorf("parse clash yaml failed: %w", err)
	}
	// 提取proxies
	proxies, ok := data["proxies"].([]any)
	if !ok {
		return nil, fmt.Errorf("empty proxies in clash yaml")
	}
	parseNodes := make([]types.ProxyNode, 0)
	for _, proxy := range proxies {
		item, ok := proxy.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid proxy item: %v", proxy)
		}
		for _, c := range p.converters {
			if !c.Handle(item) {
				continue
			}
			list, err := c.Parse(item)
			if err != nil {
				return nil, fmt.Errorf("parse shadowsocks failed: %w", err)
			}
			for i, node := range list {
				if node.Type == "" {
					list[i].Type = c.NodeType()
				}
				list[i].FromSub = sub
			}
			parseNodes = append(parseNodes, list...)
			break
		}
	}
	return parseNodes, nil
}

func (p Clash) Convert2SingBoxOutbounds(node types.ProxyNode) map[string]any {
	return p.converterMap[node.Type].Convert2SingBox(node)
}
