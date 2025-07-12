package converter

import (
	clash2 "sing-box-sub-converter/converter/clash"
	content2 "sing-box-sub-converter/converter/content"
	"sing-box-sub-converter/converter/types"
)

type SubscriptionParser interface {
	// SubType 订阅类型
	SubType() string
	// Parse 从订阅数据中解析节点列表
	Parse(content, sub string) ([]types.ProxyNode, error)
	// Convert2SingBoxOutbounds 将节点列表转换为SingBox格式
	Convert2SingBoxOutbounds(types.ProxyNode) map[string]any
}

var (
	clash = clash2.NewClash()

	content = content2.NewContent()
)

func parsers() []SubscriptionParser {
	return []SubscriptionParser{
		clash,
		content,
	}
}

func GetParser(subType string) SubscriptionParser {
	switch subType {
	case clash.SubType():
		return clash
	case content.SubType():
		return content
	}
	return nil
}
