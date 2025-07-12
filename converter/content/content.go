package content

import (
	"fmt"
	"log/slog"
	"sing-box-sub-converter/converter/types"
	"sing-box-sub-converter/utils"
	"strings"
)

type ProtocolConverter interface {
	NodeType() types.ProxyNodeType
	Handle(string) bool
	Parse(string) ([]types.ProxyNode, error)
	Convert2SingBox(types.ProxyNode) map[string]any
}

type Content struct {
	converters   []ProtocolConverter
	converterMap map[types.ProxyNodeType]ProtocolConverter
}

func NewContent() *Content {
	converters := make([]ProtocolConverter, 0)

	converters = append(converters, shadowsocks{})

	converterMap := make(map[types.ProxyNodeType]ProtocolConverter)
	for _, protocolConverter := range converters {
		converterMap[protocolConverter.NodeType()] = protocolConverter
	}
	return &Content{
		converters,
		converterMap,
	}
}

func (Content) SubType() string {
	return "content"
}

func (p Content) Parse(content, sub string) ([]types.ProxyNode, error) {
	content = strings.TrimSpace(content)
	if !strings.Contains(content, "\n") {
		//只有一行，姑且认为是base64了
		r, err := utils.Base64Decode(content)
		if err == nil {
			content = r
		} else {
			slog.Info("try to base64 decode failed")
		}
	}
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty content")
	}
	proxies := make([]string, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		proxies = append(proxies, line)
	}
	parseNodes := make([]types.ProxyNode, 0)
	for _, proxy := range proxies {
		for _, c := range p.converters {
			if !c.Handle(proxy) {
				continue
			}
			list, err := c.Parse(proxy)
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

func (p Content) Convert2SingBoxOutbounds(node types.ProxyNode) map[string]any {
	return p.converterMap[node.Type].Convert2SingBox(node)
}
