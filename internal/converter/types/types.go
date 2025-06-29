package types

type ProxyNodeType string

const (
	ProxyNodeTypeShadowsocks ProxyNodeType = "shadowsocks"
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

type ShadowsocksNode struct {
	Server     string
	ServerPort int
	Method     string
	Password   string
	Plugin     string
	PluginOpts string
}
