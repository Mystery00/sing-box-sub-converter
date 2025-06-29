package clash

import (
	"fmt"
	"sing-box-sub-converter/internal/converter/types"
	"sing-box-sub-converter/utils"
	"strconv"
	"strings"
)

type shadowsocks struct {
}

func (shadowsocks) NodeType() types.ProxyNodeType {
	return types.ProxyNodeTypeShadowsocks
}

func (shadowsocks) Handle(m map[string]any) bool {
	if d, exist := m["type"]; exist {
		return d.(string) == "ss"
	}
	return false
}

func (s shadowsocks) Parse(m map[string]any) (resultList []types.ProxyNode, err error) {
	resultList = make([]types.ProxyNode, 0)
	innerNode := types.ShadowsocksNode{}
	n := types.ProxyNode{}
	if d, exist := m["name"]; exist {
		n.Tag = d.(string)
	} else {
		n.Tag = fmt.Sprintf("%s_shadowsocks", utils.GenName(8))
	}
	if d, exist := m["server"]; exist {
		n.Address = d.(string)
		innerNode.Server = d.(string)
	}
	if d, exist := m["port"]; exist {
		p := d.(int)
		n.Port = fmt.Sprintf("%d", p)
		innerNode.ServerPort = p
	}
	if d, exist := m["cipher"]; exist {
		innerNode.Method = d.(string)
	}
	if d, exist := m["password"]; exist {
		innerNode.Password = d.(string)
	}
	plugin, pluginExist := m["plugin"]
	pluginOpts, pluginOptsExist := m["plugin-opts"]
	opts := make([]string, 0)
	if pluginExist && pluginOptsExist {
		pluginStr := plugin.(string)
		pluginOptsMap := pluginOpts.(map[string]any)
		if pluginStr == "obfs-local" || pluginStr == "obfs" {
			opts = append(opts, fmt.Sprintf("obfs=%s", pluginOptsMap["mode"]))
			opts = append(opts, fmt.Sprintf("obfs-host=%s", pluginOptsMap["host"]))
			innerNode.Plugin = "obfs-local"
		} else if pluginStr == "v2ray-plugin" {
			opts = append(opts, fmt.Sprintf("obfs=%s", pluginOptsMap["mode"]))
			opts = append(opts, fmt.Sprintf("obfs-host=%s", pluginOptsMap["host"]))
			if d, exist := pluginOptsMap["path"]; exist {
				opts = append(opts, fmt.Sprintf("path=%s", d))
			}
			if d, exist := pluginOptsMap["headers"]; exist {
				opts = append(opts, fmt.Sprintf("headers=%s", utils.JsonStr(d)))
			}
			if d, exist := pluginOptsMap["fingerprint"]; exist {
				opts = append(opts, fmt.Sprintf("fingerprint=%s", d))
			}
			if d, exist := pluginOptsMap["mux"]; exist {
				opts = append(opts, fmt.Sprintf("mux=%s", d))
			}
			if d, exist := pluginOptsMap["skip-cert-verify"]; exist {
				opts = append(opts, fmt.Sprintf("skip-cert-verify=%s", d))
			}
			if d, exist := pluginOptsMap["tls"]; exist {
				opts = append(opts, fmt.Sprintf("tls=%s", d))
			}
			innerNode.Plugin = "v2ray-plugin"
		}
	}
	innerNode.PluginOpts = strings.Join(opts, ";")

	if innerNode.Method == "chacha20-poly1305" {
		innerNode.Method = "chacha20-ietf-poly1305"
	}
	if innerNode.Method == "xchacha20-poly1305" {
		innerNode.Method = "xchacha20-ietf-poly1305"
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
	innerNode := node.ProxyDetail.(types.ShadowsocksNode)
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
