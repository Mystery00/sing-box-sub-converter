package clash

import (
	"fmt"
	"sing-box-sub-converter/converter/types"
	"sing-box-sub-converter/utils"
	"strconv"
	"strings"
)

type VmessNode struct {
	Uuid                string
	AlterId             int
	Security            string
	GlobalPadding       bool
	AuthenticatedLength bool
	Tls                 map[string]any
	Transport           map[string]any
}

type vmess struct {
}

func (vmess) NodeType() types.ProxyNodeType {
	return types.ProxyNodeTypeVmess
}

func (vmess) Handle(m map[string]any) bool {
	if d, exist := m["type"]; exist {
		return d.(string) == "vmess"
	}
	return false
}

func (v vmess) Parse(m map[string]any) (resultList []types.ProxyNode, err error) {
	resultList = make([]types.ProxyNode, 0)
	innerNode := VmessNode{
		Tls:       make(map[string]any),
		Transport: make(map[string]any),
	}
	n := types.ProxyNode{}

	// 解析节点名称
	if d, exist := m["name"]; exist {
		n.Tag = strings.TrimSpace(d.(string))
	} else {
		n.Tag = fmt.Sprintf("%s_vmess", utils.GenName(8))
	}

	// 解析服务器地址
	if d, exist := m["server"]; exist {
		n.Address = strings.TrimSpace(d.(string))
	}

	// 解析端口
	if d, exist := m["port"]; exist {
		p := d.(int)
		n.Port = fmt.Sprintf("%d", p)
	}

	// 解析UUID
	if d, exist := m["uuid"]; exist {
		innerNode.Uuid = d.(string)
	}

	// 解析alterId
	if d, exist := m["alterId"]; exist {
		innerNode.AlterId = d.(int)
	}

	// 解析加密方式 (cipher)
	if d, exist := m["cipher"]; exist {
		security := d.(string)
		if security == "auto" {
			security = "auto"
		}
		innerNode.Security = security
	} else {
		innerNode.Security = "auto"
	}

	// 解析TLS相关配置
	if d, exist := m["tls"]; exist && d.(bool) {
		innerNode.Tls["enabled"] = true
	}

	if d, exist := m["skip-cert-verify"]; exist {
		if len(innerNode.Tls) == 0 {
			innerNode.Tls = make(map[string]any)
		}
		innerNode.Tls["insecure"] = d.(bool)
	}

	if d, exist := m["sni"]; exist {
		if len(innerNode.Tls) == 0 {
			innerNode.Tls = make(map[string]any)
		}
		innerNode.Tls["server_name"] = d.(string)
	}

	if d, exist := m["alpn"]; exist {
		if len(innerNode.Tls) == 0 {
			innerNode.Tls = make(map[string]any)
		}
		innerNode.Tls["alpn"] = d
	}

	// 解析传输协议 (network)
	if d, exist := m["network"]; exist {
		network := d.(string)
		switch network {
		case "ws":
			innerNode.Transport["type"] = "ws"
			if d, exist := m["ws-opts"]; exist {
				wsOpts := d.(map[string]any)
				if path, exist := wsOpts["path"]; exist {
					innerNode.Transport["path"] = path
				}
				if headers, exist := wsOpts["headers"]; exist {
					innerNode.Transport["headers"] = headers
				}
			}
		case "h2":
			innerNode.Transport["type"] = "http"
			if d, exist := m["h2-opts"]; exist {
				h2Opts := d.(map[string]any)
				if path, exist := h2Opts["path"]; exist {
					innerNode.Transport["path"] = path
				}
				if host, exist := h2Opts["host"]; exist {
					innerNode.Transport["host"] = host
				}
			}
		case "grpc":
			innerNode.Transport["type"] = "grpc"
			if d, exist := m["grpc-opts"]; exist {
				grpcOpts := d.(map[string]any)
				if serviceName, exist := grpcOpts["grpc-service-name"]; exist {
					innerNode.Transport["service_name"] = serviceName
				}
			}
		}
	}

	// 如果没有配置TLS，则设置为nil
	if len(innerNode.Tls) == 0 {
		innerNode.Tls = nil
	}

	// 如果没有配置传输协议，则设置为nil
	if len(innerNode.Transport) == 0 {
		innerNode.Transport = nil
	}

	n.ProxyDetail = innerNode
	resultList = append(resultList, n)
	return resultList, nil
}

func (vmess) Convert2SingBox(node types.ProxyNode) map[string]any {
	innerNode := node.ProxyDetail.(VmessNode)
	m := make(map[string]any)

	m["tag"] = node.Tag
	m["type"] = node.Type
	m["server"] = node.Address
	pp, _ := strconv.Atoi(node.Port)
	m["server_port"] = uint16(pp)
	m["uuid"] = innerNode.Uuid
	m["security"] = innerNode.Security

	if innerNode.AlterId != 0 {
		m["alter_id"] = innerNode.AlterId
	}

	if innerNode.GlobalPadding {
		m["global_padding"] = innerNode.GlobalPadding
	}

	if innerNode.AuthenticatedLength {
		m["authenticated_length"] = innerNode.AuthenticatedLength
	}

	if innerNode.Tls != nil {
		m["tls"] = innerNode.Tls
	}

	if innerNode.Transport != nil {
		m["transport"] = innerNode.Transport
	}

	return m
}
