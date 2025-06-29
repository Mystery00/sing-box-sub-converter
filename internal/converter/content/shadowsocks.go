package content

import (
	"fmt"
	"net/url"
	"sing-box-sub-converter/internal/converter/types"
	"sing-box-sub-converter/utils"
	"strconv"
	"strings"
)

type ShadowsocksNode struct {
	Method     string
	Password   string
	Plugin     string
	PluginOpts string
}

type shadowsocks struct {
}

func (shadowsocks) NodeType() types.ProxyNodeType {
	return types.ProxyNodeTypeShadowsocks
}

func (shadowsocks) Handle(content string) bool {
	return strings.HasPrefix(content, "ss://")
}

func (s shadowsocks) Parse(content string) (resultList []types.ProxyNode, err error) {
	resultList = make([]types.ProxyNode, 0)
	u, err := url.Parse(content)
	if err != nil {
		return resultList, fmt.Errorf("parse uri failed: %w", err)
	}
	innerNode := ShadowsocksNode{}
	n := types.ProxyNode{}
	if u.Fragment != "" {
		n.Tag = strings.TrimSpace(u.Fragment)
	} else {
		n.Tag = fmt.Sprintf("%s_shadowsocks", utils.GenName(8))
	}
	n.Address = u.Hostname()
	n.Port = u.Port()
	decode, err := utils.Base64UrlDecode(u.User.Username())
	if err != nil {
		return resultList, fmt.Errorf("base64 decode failed: %w", err)
	}
	index := strings.Index(decode, ":")
	innerNode.Method = decode[0:index]
	innerNode.Password = decode[index+1:]
	if u.Query().Get("plugin") != "" {
		pluginStr := u.Query().Get("plugin")
		i := strings.Index(pluginStr, ";")
		innerNode.Plugin = pluginStr[0:i]
		innerNode.PluginOpts = pluginStr[i+1:]
	}
	n.ProxyDetail = innerNode
	resultList = append(resultList, n)
	return resultList, nil
}

func (shadowsocks) Convert2SingBox(node types.ProxyNode) map[string]any {
	m := make(map[string]any)
	m["tag"] = node.Tag
	m["type"] = node.Type
	m["server"] = node.Address
	pp, _ := strconv.Atoi(node.Port)
	m["server_port"] = uint16(pp)
	innerNode := node.ProxyDetail.(ShadowsocksNode)
	m["method"] = innerNode.Method
	m["password"] = innerNode.Password
	if innerNode.Plugin != "" {
		m["plugin"] = innerNode.Plugin
	}
	if innerNode.PluginOpts != "" {
		m["plugin_opts"] = innerNode.PluginOpts
	}
	return m
}
