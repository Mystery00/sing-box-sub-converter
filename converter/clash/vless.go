package clash

import (
	"fmt"
	"sing-box-sub-converter/converter/types"
	"sing-box-sub-converter/utils"
	"strconv"
	"strings"
)

type Vless2Node struct {
	Uuid           string
	Flow           string
	PacketEncoding string
	Tls            map[string]any
	Reality        map[string]any
}

type vless struct {
}

func (vless) NodeType() types.ProxyNodeType {
	return types.ProxyNodeTypeVless
}

func (vless) Handle(m map[string]any) bool {
	if d, exist := m["type"]; exist {
		return d.(string) == "vless"
	}
	return false
}

func (v vless) Parse(m map[string]any) (resultList []types.ProxyNode, err error) {
	resultList = make([]types.ProxyNode, 0)
	innerNode := Vless2Node{
		Tls:     make(map[string]any),
		Reality: make(map[string]any),
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
	if d, exist := m["uuid"]; exist {
		innerNode.Uuid = d.(string)
	}
	if d, exist := m["flow"]; exist {
		innerNode.Flow = d.(string)
	}
	if d, exist := m["skip-cert-verify"]; exist {
		innerNode.Tls["insecure"] = d.(bool)
	}
	if d, exist := m["sni"]; exist {
		innerNode.Tls["server_name"] = d.(string)
	}
	if d, exist := m["reality-opts"]; exist {
		mm := d.(map[string]any)
		innerNode.Reality["security"] = true
		if d, exist := mm["public-key"]; exist {
			innerNode.Reality["public_key"] = d.(string)
		}
		if d, exist := mm["short-id"]; exist {
			innerNode.Reality["short_id"] = d.(string)
		}
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

func (vless) Convert2SingBox(node types.ProxyNode) map[string]any {
	innerNode := node.ProxyDetail.(Vless2Node)
	m := make(map[string]any)
	m["tag"] = node.Tag
	m["type"] = node.Type
	m["server"] = node.Address
	pp, _ := strconv.Atoi(node.Port)
	m["server_port"] = uint16(pp)
	m["uuid"] = innerNode.Uuid
	m["flow"] = innerNode.Flow
	if innerNode.Tls != nil {
		m["tls"] = innerNode.Tls
	}
	return m
}
