package template

import (
	"encoding/json"
	"fmt"
	"github.com/deckarep/golang-set/v2"
	"log/slog"
	"os"
	"path/filepath"
	"sing-box-sub-converter/converter"
	"sing-box-sub-converter/converter/types"
	"strings"
)

// GetConfigTemplate 获取配置模板
func GetConfigTemplate(templateName string) (map[string]any, error) {
	var config map[string]any

	// 读取环境变量中的模板目录路径,默认为config_template
	templateDir := os.Getenv("TEMPLATE_DIR")
	if templateDir == "" {
		templateDir = "config_templates"
	}

	// 从本地文件加载模板
	templateList, err := loadTemplateList(templateDir)
	if err != nil {
		return nil, fmt.Errorf("获取模板列表失败: %w", err)
	}

	if len(templateList) < 1 {
		return nil, fmt.Errorf("没有找到模板文件")
	}

	// 加载选定的模板
	templateFileName := fmt.Sprintf("%s.json", templateName)
	templateFilePath := filepath.Join(templateDir, templateFileName)

	templateData, err := os.ReadFile(templateFilePath)
	if err != nil {
		return nil, fmt.Errorf("读取模板文件失败: %w", err)
	}

	if err := json.Unmarshal(templateData, &config); err != nil {
		return nil, fmt.Errorf("解析模板文件失败: %w", err)
	}

	return config, nil
}

// MergeToConfig 将节点合并到配置模板中
func MergeToConfig(config map[string]any, nodes []types.ProxyNode) (map[string]any, error) {
	// 获取outbounds
	outbounds, ok := config["outbounds"].([]any)
	if !ok {
		return nil, fmt.Errorf("无效的outbounds配置")
	}
	validNodes := make([]types.ProxyNode, 0)
	subInfoNodes := make([]string, 0)
	for _, node := range nodes {
		if node.Type == types.ProxyNodeTypeSubInfo {
			subInfoNodes = append(subInfoNodes, node.Tag)
			continue
		}
		validNodes = append(validNodes, node)
	}

	// 依次替换规则中的占位符
	usedNodeTags := mapset.NewSet[string]()
	for i, o := range outbounds {
		outbound, ok := o.(map[string]any)
		if !ok {
			continue
		}
		if outbound["outbounds"] == nil {
			continue
		}
		originOutbounds, ok := outbound["outbounds"].([]any)
		if !ok {
			continue
		}
		replaceNodes := make([]types.ProxyNode, len(validNodes))
		copy(replaceNodes, validNodes)
		filterVal, exist := outbound["filter"]
		if exist {
			filters := parseFilters(filterVal)
			replaceNodes = handleFilter(replaceNodes, filters)
		}
		delete(outbound, "filter")
		resultOutbounds := make([]string, 0)
		for _, oo := range originOutbounds {
			if oo == "{all}" {
				//填充所有节点
				for _, node := range replaceNodes {
					resultOutbounds = append(resultOutbounds, node.Tag)
					usedNodeTags.Add(node.Tag)
				}
				continue
			}
			ooStr := fmt.Sprintf("%v", oo)
			if strings.HasPrefix(ooStr, "{") && strings.HasSuffix(ooStr, "}") {
				nodeTag := ooStr[1 : len(ooStr)-1]
				for _, node := range replaceNodes {
					if node.FromSub == nodeTag {
						resultOutbounds = append(resultOutbounds, node.Tag)
						usedNodeTags.Add(node.Tag)
					}
				}
				continue
			}
			resultOutbounds = append(resultOutbounds, ooStr)
		}
		if len(resultOutbounds) == 0 {
			slog.Warn(fmt.Sprintf("发现 %v 出站下的节点数量为 0 ，会导致sing-box无法运行，请检查config模板是否正确", outbound["tag"]))
		}
		outbound["outbounds"] = resultOutbounds
		outbounds[i] = outbound
	}
	if len(subInfoNodes) != 0 {
		//这个方法中添加一个选择器
		outbound := make(map[string]any)
		outbound["type"] = "selector"
		outbound["tag"] = "订阅信息"
		outbound["outbounds"] = subInfoNodes
		outbounds = append(outbounds, outbound)
	}

	//将所有节点加到最后
	appendNodes := make([]types.ProxyNode, 0)
	for _, node := range nodes {
		if usedNodeTags.Contains(node.Tag) {
			appendNodes = append(appendNodes, node)
		}
	}

	for _, node := range appendNodes {
		o := converter.GetParser(node.SubType).Convert2SingBoxOutbounds(node)
		outbounds = append(outbounds, o)
	}
	if len(subInfoNodes) != 0 {
		//把订阅信息单独加到出站列表
		for _, tag := range subInfoNodes {
			outbound := make(map[string]any)
			outbound["type"] = "direct"
			outbound["tag"] = tag
			outbounds = append(outbounds, outbound)
		}
	}

	config["outbounds"] = outbounds
	return config, nil
}
