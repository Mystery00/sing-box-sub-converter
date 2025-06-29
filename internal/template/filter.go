package template

import (
	"fmt"
	"sing-box-sub-converter/internal/converter/types"
	"strings"
)

func handleFilter(nodes []types.ProxyNode, filters []OutboundFilter) []types.ProxyNode {
	var results []types.ProxyNode
	for _, filter := range filters {
		if len(filter.Keywords) == 0 {
			continue
		}
		switch filter.Action {
		case "include":
			results = filterIncludeKeywords(nodes, filter.Keywords, filter.ForTag)
			break
		case "exclude":
			results = filterExcludeKeywords(nodes, filter.Keywords, filter.ForTag)
			break
		default:
			break
		}
	}
	return results
}

func filterIncludeKeywords(nodes []types.ProxyNode, keywords []string, forTag string) []types.ProxyNode {
	results := make([]types.ProxyNode, 0)
	for _, result := range nodes {
		if forTag != "" && forTag != result.FromSub {
			results = append(results, result)
			continue
		}
		for _, keyword := range keywords {
			if strings.Contains(result.Tag, keyword) {
				results = append(results, result)
				break
			}
		}
	}
	return results
}

func filterExcludeKeywords(nodes []types.ProxyNode, keywords []string, forTag string) []types.ProxyNode {
	results := make([]types.ProxyNode, 0)
	for _, result := range nodes {
		if forTag != "" && forTag != result.FromSub {
			results = append(results, result)
			continue
		}
		shouldBreak := false
		for _, keyword := range keywords {
			if strings.Contains(result.Tag, keyword) {
				shouldBreak = true
				break
			}
		}
		if shouldBreak {
			continue
		}
		results = append(results, result)
	}
	return results
}

func parseFilters(content any) []OutboundFilter {
	filters := make([]OutboundFilter, 0)
	list, ok := content.([]any)
	if !ok {
		return filters
	}
	for _, f := range list {
		ff, ok := f.(map[string]any)
		if !ok {
			continue
		}
		filter := OutboundFilter{}
		if d, ok := ff["action"].(string); ok {
			filter.Action = d
		}
		if d, ok := ff["keywords"].([]any); ok {
			keywords := make([]string, 0)
			for _, s := range d {
				ss := strings.TrimSpace(fmt.Sprintf("%v", s))
				if ss != "" {
					keywords = append(keywords, ss)
				}
			}
			filter.Keywords = keywords
		}
		if d, ok := ff["for"].(string); ok {
			filter.ForTag = d
		}
		if filter.Action == "" && len(filter.Keywords) == 0 && filter.ForTag == "" {
			continue
		}
		filters = append(filters, filter)
	}
	return filters
}
