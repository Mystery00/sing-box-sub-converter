package types

type ProxyNodeType string

const (
	ProxyNodeTypeShadowsocks ProxyNodeType = "shadowsocks"
	ProxyNodeTypeTrojan      ProxyNodeType = "trojan"
	ProxyNodeTypeHysteria2   ProxyNodeType = "hysteria2"
	ProxyNodeTypeVless       ProxyNodeType = "vless"
	ProxyNodeTypeVmess       ProxyNodeType = "vmess"
)

type ProxyNode struct {
	Type        ProxyNodeType
	Tag         string
	Address     string
	Port        string
	FromSub     string
	SubType     string
	ProxyDetail any
}
