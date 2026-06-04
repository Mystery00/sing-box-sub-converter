package clash

import (
	"fmt"
	"sing-box-sub-converter/converter/types"
	"sing-box-sub-converter/utils"
	"strconv"
)

type AnyTlsNode struct {
	Password                 string
	IdleSessionCheckInterval string
	IdleSessionTimeout       string
	MinIdleSession           int
	Detour                   string
	Tls                      map[string]any
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
	innerNode := AnyTlsNode{
		Tls: map[string]any{
			"enabled":  true,
			"insecure": false,
		},
	}
	if d, exist := getString(m, "dialer-proxy"); exist {
		innerNode.Detour = d
	} else if d, exist := getString(m, "detour"); exist {
		innerNode.Detour = d
	}

	n := types.ProxyNode{}
	if d, exist := getString(m, "name"); exist {
		n.Tag = d
	} else {
		n.Tag = fmt.Sprintf("%s_anytls", utils.GenName(8))
	}
	if d, exist := getString(m, "server"); exist {
		n.Address = d
		innerNode.Tls["server_name"] = d
	}
	if p, exist := getInt(m, "port"); exist {
		n.Port = fmt.Sprintf("%d", p)
	}
	if d, exist := getString(m, "password"); exist {
		innerNode.Password = d
	}
	if p, exist := getInt(m, "idle-session-check-interval"); exist {
		innerNode.IdleSessionCheckInterval = fmt.Sprintf("%ds", p)
	}
	if p, exist := getInt(m, "idle-session-timeout"); exist {
		innerNode.IdleSessionTimeout = fmt.Sprintf("%ds", p)
	}
	if d, exist := getInt(m, "min-idle-session"); exist {
		innerNode.MinIdleSession = d
	}
	if d, exist := getBool(m, "skip-cert-verify"); exist {
		innerNode.Tls["insecure"] = d
	}
	if d, exist := getString(m, "sni"); exist {
		innerNode.Tls["server_name"] = d
	}
	if d, exist := getStringSlice(m, "alpn"); exist {
		innerNode.Tls["alpn"] = d
	}
	if d, exist := getString(m, "client-fingerprint"); exist {
		if fingerprint, ok := normalizeUtlsFingerprint(d); ok {
			innerNode.Tls["utls"] = map[string]any{
				"enabled":     true,
				"fingerprint": fingerprint,
			}
		}
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
	if innerNode.IdleSessionCheckInterval != "" {
		m["idle_session_check_interval"] = innerNode.IdleSessionCheckInterval
	}
	if innerNode.IdleSessionTimeout != "" {
		m["idle_session_timeout"] = innerNode.IdleSessionTimeout
	}
	if innerNode.MinIdleSession != 0 {
		m["min_idle_session"] = innerNode.MinIdleSession
	}
	if innerNode.Detour != "" {
		m["detour"] = innerNode.Detour
	}
	if innerNode.Tls != nil {
		m["tls"] = innerNode.Tls
	}
	return m
}
