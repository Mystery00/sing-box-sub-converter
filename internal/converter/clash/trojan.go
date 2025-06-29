package clash

import (
	"fmt"
	"sing-box-sub-converter/internal/converter/types"
	"sing-box-sub-converter/utils"
	"strconv"
	"strings"
)

type TrojanNode struct {
	Password string
	Network  string
	Tls      map[string]any
}

type trojan struct {
}

func (trojan) NodeType() types.ProxyNodeType {
	return types.ProxyNodeTypeTrojan
}

func (trojan) Handle(m map[string]any) bool {
	if d, exist := m["type"]; exist {
		return d.(string) == "trojan"
	}
	return false
}

func (t trojan) Parse(m map[string]any) (resultList []types.ProxyNode, err error) {
	resultList = make([]types.ProxyNode, 0)
	innerNode := TrojanNode{
		Tls: make(map[string]any),
	}
	n := types.ProxyNode{}
	if d, exist := m["name"]; exist {
		n.Tag = strings.TrimSpace(d.(string))
	} else {
		n.Tag = fmt.Sprintf("%s_trojan", utils.GenName(8))
	}
	if d, exist := m["server"]; exist {
		n.Address = strings.TrimSpace(d.(string))
	}
	if d, exist := m["port"]; exist {
		p := d.(int)
		n.Port = fmt.Sprintf("%d", p)
	}
	if d, exist := m["password"]; exist {
		innerNode.Password = d.(string)
	}
	if d, exist := m["skip-cert-verify"]; exist {
		innerNode.Tls["insecure"] = d.(bool)
	}
	if d, exist := m["alpn"]; exist {
		innerNode.Tls["alpn"] = d
	}
	if d, exist := m["sni"]; exist {
		innerNode.Tls["server_name"] = d.(string)
	}
	if d, exist := m["client-fingerprint"]; exist {
		innerNode.Tls["utls"] = map[string]any{
			"enabled":     true,
			"fingerprint": d.(string),
		}
	}
	if len(innerNode.Tls) == 0 {
		innerNode.Tls = nil
	} else {
		innerNode.Tls["enabled"] = true
	}

	n.ProxyDetail = innerNode
	resultList = append(resultList, n)
	return resultList, nil
}

func (trojan) Convert2SingBox(node types.ProxyNode) map[string]any {
	innerNode := node.ProxyDetail.(TrojanNode)
	m := make(map[string]any)
	m["tag"] = node.Tag
	m["type"] = node.Type
	m["server"] = node.Address
	pp, _ := strconv.Atoi(node.Port)
	m["server_port"] = uint16(pp)
	m["password"] = innerNode.Password
	if innerNode.Tls != nil {
		m["tls"] = innerNode.Tls
	}
	return m
}
