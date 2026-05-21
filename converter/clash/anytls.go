package clash

import (
	"fmt"
	"sing-box-sub-converter/converter/types"
	"sing-box-sub-converter/utils"
	"strconv"
	"strings"
)

type AnyTlsNode struct {
	Password                 string
	IdleSessionCheckInterval string
	IdleSessionTimeout       string
	MinIdleSession           int
	TcpFastOpen              bool
	Tls                      map[string]any
	UTls                     map[string]any
}

type anytls struct {
}

func (anytls) NodeType() types.ProxyNodeType {
	return types.ProxyNodeTypeAnyTls
}

func (anytls) Handle(m map[string]any) bool {
	if d, exist := m["type"]; exist {
		return d.(string) == "anytls"
	}
	return false
}

func (s anytls) Parse(m map[string]any) (resultList []types.ProxyNode, err error) {
	resultList = make([]types.ProxyNode, 0)
	innerNode := AnyTlsNode{}
	n := types.ProxyNode{}
	if d, exist := m["name"]; exist {
		n.Tag = strings.TrimSpace(d.(string))
	} else {
		n.Tag = fmt.Sprintf("%s_anytls", utils.GenName(8))
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
	if d, exist := m["idle-session-check-interval"]; exist {
		p := d.(int)
		innerNode.IdleSessionCheckInterval = fmt.Sprintf("%ds", p)
	}
	if d, exist := m["idle-session-timeout"]; exist {
		p := d.(int)
		innerNode.IdleSessionTimeout = fmt.Sprintf("%ds", p)
	}
	if d, exist := m["min-idle-session"]; exist {
		innerNode.MinIdleSession = d.(int)
	}
	if d, exist := m["tfo"]; exist {
		innerNode.TcpFastOpen = d.(bool)
	}

	innerNode.Tls = make(map[string]any)
	innerNode.UTls = make(map[string]any)

	// 解析TLS相关配置
	innerNode.Tls["enabled"] = true

	if d, exist := m["skip-cert-verify"]; exist {
		innerNode.Tls["insecure"] = d.(bool)
	}

	if d, exist := m["sni"]; exist {
		if len(innerNode.Tls) == 0 {
			innerNode.Tls = make(map[string]any)
		}
		innerNode.Tls["server_name"] = d.(string)
	}

	if d, exist := m["alpn"]; exist {
		innerNode.Tls["alpn"] = d
	}

	if d, exist := m["client-fingerprint"]; exist {
		innerNode.UTls["enabled"] = true
		innerNode.UTls["fingerprint"] = d.(string)
	}

	// 如果没有配置TLS，则设置为nil
	if len(innerNode.Tls) == 0 {
		innerNode.Tls = nil
	}
	if len(innerNode.UTls) == 0 {
		innerNode.UTls = nil
	}

	n.ProxyDetail = innerNode
	resultList = append(resultList, n)
	return resultList, nil
}

func (anytls) Convert2SingBox(node types.ProxyNode) map[string]any {
	innerNode := node.ProxyDetail.(AnyTlsNode)
	m := make(map[string]any)
	m["tag"] = node.Tag
	m["type"] = node.Type
	m["server"] = node.Address
	pp, _ := strconv.Atoi(node.Port)
	m["server_port"] = uint16(pp)
	m["password"] = innerNode.Password
	m["idle_session_check_interval"] = innerNode.IdleSessionCheckInterval
	m["idle_session_timeout"] = innerNode.IdleSessionTimeout
	m["min_idle_session"] = innerNode.MinIdleSession
	m["tcp_fast_open"] = innerNode.TcpFastOpen
	tls := innerNode.Tls
	if innerNode.UTls != nil {
		if innerNode.Tls == nil {
			innerNode.Tls = make(map[string]any)
		}
		tls["utls"] = innerNode.UTls
	}
	m["tls"] = tls
	return m
}
