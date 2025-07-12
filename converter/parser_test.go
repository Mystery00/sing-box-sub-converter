package converter

import (
	"sing-box-sub-converter/converter/types"
	"testing"
)

func Test_processDuplicateNodeTag(t *testing.T) {
	nodes := make([]types.ProxyNode, 0)
	nodes = append(nodes, types.ProxyNode{Tag: "美国1"})
	nodes = append(nodes, types.ProxyNode{Tag: "美国2"})
	nodes = append(nodes, types.ProxyNode{Tag: "美国2"})
	nodes = append(nodes, types.ProxyNode{Tag: "美国2"})
	nodes = append(nodes, types.ProxyNode{Tag: "美国2"})
	nodes = append(nodes, types.ProxyNode{Tag: "美国2_1"})
	nodes = append(nodes, types.ProxyNode{Tag: "美国2_1"})
	nodes = append(nodes, types.ProxyNode{Tag: "美国2_1"})
	nodes = append(nodes, types.ProxyNode{Tag: "美国2_1_1"})
	nodes = append(nodes, types.ProxyNode{Tag: "美国3"})
	nodes = append(nodes, types.ProxyNode{Tag: "美国4"})
	nodes = append(nodes, types.ProxyNode{Tag: "美国5"})
	nodes = append(nodes, types.ProxyNode{Tag: "美国6"})
	nodes = append(nodes, types.ProxyNode{Tag: "美国6"})
	tag := processDuplicateNodeTag(nodes)
	t.Log(tag)
}
