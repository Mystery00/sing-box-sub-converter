package template

import (
	"fmt"
	"os"
	"strings"
)

// loadTemplateList 获取配置模板列表
func loadTemplateList() ([]string, error) {
	// 配置模板文件夹路径
	templateDir := "config_template"

	// 读取文件夹内容
	files, err := os.ReadDir(templateDir)
	if err != nil {
		return nil, fmt.Errorf("读取模板目录失败: %w", err)
	}

	// 提取JSON文件名（不含扩展名）
	templateList := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			name := strings.TrimSuffix(file.Name(), ".json")
			templateList = append(templateList, name)
		}
	}
	return templateList, nil
}
