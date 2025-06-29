package clash

import (
	"fmt"
	"sing-box-sub-converter/internal/converter/types"
	"sing-box-sub-converter/utils"
	"strconv"
	"strings"
)

type Hysteria2Node struct {
	ServerPorts []string
	Password    string
	Obfs        map[string]string
	UpMbps      int
	DownMbps    int
	Tls         map[string]any
}

type hysteria2 struct {
}

func (hysteria2) NodeType() types.ProxyNodeType {
	return types.ProxyNodeTypeHysteria2
}

func (hysteria2) Handle(m map[string]any) bool {
	if d, exist := m["type"]; exist {
		return d.(string) == "hysteria2"
	}
	return false
}

func (h hysteria2) Parse(m map[string]any) (resultList []types.ProxyNode, err error) {
	resultList = make([]types.ProxyNode, 0)
	innerNode := Hysteria2Node{
		Tls: make(map[string]any),
	}
	n := types.ProxyNode{}
	if d, exist := m["name"]; exist {
		n.Tag = strings.TrimSpace(d.(string))
	} else {
		n.Tag = fmt.Sprintf("%s_hysteria2", utils.GenName(8))
	}
	if d, exist := m["server"]; exist {
		n.Address = strings.TrimSpace(d.(string))
	}
	if d, exist := m["port"]; exist {
		p := d.(int)
		n.Port = fmt.Sprintf("%d", p)
	}
	if d, exist := m["ports"]; exist {
		p := d.([]any)
		innerNode.ServerPorts = make([]string, 0)
		for _, v := range p {
			vv := fmt.Sprintf("%v", v)
			if vv == "" {
				continue
			}
			innerNode.ServerPorts = append(innerNode.ServerPorts, vv)
		}
	}
	if d, exist := m["password"]; exist {
		innerNode.Password = d.(string)
	}
	if d, exist := m["skip-cert-verify"]; exist {
		innerNode.Tls["insecure"] = d.(bool)
	}
	if d, exist := m["obfs"]; exist {
		innerNode.Obfs["type"] = d.(string)
	}
	if d, exist := m["obfs-password"]; exist {
		innerNode.Obfs["password"] = d.(string)
	}
	if d, exist := m["fingerprint"]; exist {
		innerNode.Obfs["password"] = d.(string)
	}
	if d, exist := m["alpn"]; exist {
		innerNode.Tls["alpn"] = d
	}
	if d, exist := m["sni"]; exist {
		innerNode.Tls["server_name"] = d.(string)
	}
	if d, exist := m["up"]; exist {
		innerNode.UpMbps = d.(int)
	}
	if d, exist := m["down"]; exist {
		innerNode.DownMbps = d.(int)
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

func (hysteria2) Convert2SingBox(node types.ProxyNode) map[string]any {
	innerNode := node.ProxyDetail.(Hysteria2Node)
	m := make(map[string]any)
	m["tag"] = node.Tag
	m["type"] = node.Type
	m["server"] = node.Address
	pp, _ := strconv.Atoi(node.Port)
	m["server_port"] = uint16(pp)
	if len(innerNode.ServerPorts) != 0 {
		m["server_ports"] = innerNode.ServerPorts
	}
	if innerNode.UpMbps != 0 {
		m["up_mbps"] = uint16(innerNode.UpMbps)
	}
	if innerNode.DownMbps != 0 {
		m["down_mbps"] = uint16(innerNode.DownMbps)
	}
	if len(innerNode.Obfs) != 0 {
		m["obfs"] = innerNode.Obfs
	}
	m["password"] = innerNode.Password
	if innerNode.Tls != nil {
		m["tls"] = innerNode.Tls
	}
	return m
}
