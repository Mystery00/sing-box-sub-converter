package content

import (
	"encoding/json"
	"fmt"
	"sing-box-sub-converter/converter/types"
	"sing-box-sub-converter/utils"
	"strconv"
	"strings"
)

type vmessURIConfig struct {
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port any    `json:"port"`
	Id   string `json:"id"`
	Aid  any    `json:"aid"`
	Scy  string `json:"scy"`
	Net  string `json:"net"`
	Host string `json:"host"`
	Path string `json:"path"`
	Tls  string `json:"tls"`
	Sni  string `json:"sni"`
	Alpn string `json:"alpn"`
}

type VmessContentNode struct {
	Uuid      string
	AlterId   int
	Security  string
	Tls       map[string]any
	Transport map[string]any
}

type vmess struct{}

func (vmess) NodeType() types.ProxyNodeType {
	return types.ProxyNodeTypeVmess
}

func (vmess) Handle(content string) bool {
	return strings.HasPrefix(content, "vmess://")
}

func (vmess) Parse(content string) ([]types.ProxyNode, error) {
	b64 := strings.TrimPrefix(content, "vmess://")
	decoded, err := utils.Base64UrlDecode(b64)
	if err != nil {
		if decoded, err = utils.Base64Decode(b64); err != nil {
			return nil, fmt.Errorf("解码vmess URI失败: %w", err)
		}
	}
	var cfg vmessURIConfig
	if err = json.Unmarshal([]byte(decoded), &cfg); err != nil {
		return nil, fmt.Errorf("解析vmess JSON失败: %w", err)
	}

	n := types.ProxyNode{}
	if cfg.Ps != "" {
		n.Tag = strings.TrimSpace(cfg.Ps)
	} else {
		n.Tag = fmt.Sprintf("%s_vmess", utils.GenName(8))
	}
	n.Address = cfg.Add
	switch p := cfg.Port.(type) {
	case float64:
		n.Port = fmt.Sprintf("%d", int(p))
	case string:
		n.Port = p
	}

	security := cfg.Scy
	if security == "" {
		security = "auto"
	}
	alterId := 0
	switch a := cfg.Aid.(type) {
	case float64:
		alterId = int(a)
	case string:
		alterId, _ = strconv.Atoi(a)
	}

	inner := VmessContentNode{Uuid: cfg.Id, AlterId: alterId, Security: security}

	if cfg.Tls == "tls" {
		inner.Tls = map[string]any{"enabled": true}
		if cfg.Sni != "" {
			inner.Tls["server_name"] = cfg.Sni
		}
		if cfg.Alpn != "" {
			inner.Tls["alpn"] = strings.Split(cfg.Alpn, ",")
		}
	}

	switch cfg.Net {
	case "ws":
		inner.Transport = map[string]any{"type": "ws"}
		if cfg.Path != "" {
			inner.Transport["path"] = cfg.Path
		}
		if cfg.Host != "" {
			inner.Transport["headers"] = map[string]any{"Host": cfg.Host}
		}
	case "h2":
		inner.Transport = map[string]any{"type": "http"}
		if cfg.Path != "" {
			inner.Transport["path"] = cfg.Path
		}
		if cfg.Host != "" {
			inner.Transport["host"] = strings.Split(cfg.Host, ",")
		}
	case "grpc":
		inner.Transport = map[string]any{"type": "grpc"}
		if cfg.Path != "" {
			inner.Transport["service_name"] = cfg.Path
		}
	}

	n.ProxyDetail = inner
	return []types.ProxyNode{n}, nil
}

func (vmess) Convert2SingBox(node types.ProxyNode) map[string]any {
	inner := node.ProxyDetail.(VmessContentNode)
	pp, _ := strconv.Atoi(node.Port)
	m := map[string]any{
		"tag":         node.Tag,
		"type":        node.Type,
		"server":      node.Address,
		"server_port": uint16(pp),
		"uuid":        inner.Uuid,
		"security":    inner.Security,
	}
	if inner.AlterId != 0 {
		m["alter_id"] = inner.AlterId
	}
	if inner.Tls != nil {
		m["tls"] = inner.Tls
	}
	if inner.Transport != nil {
		m["transport"] = inner.Transport
	}
	return m
}
