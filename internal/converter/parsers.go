package converter

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"log/slog"
	"sing-box-sub-converter/internal/config"
	"sing-box-sub-converter/internal/converter/types"
	"sing-box-sub-converter/internal/fetcher"
	"sing-box-sub-converter/utils"
)

// ProcessSubscribes 处理所有订阅
func ProcessSubscribes(subscribes []config.Subscription) ([]types.ProxyNode, error) {
	nodes := make([]types.ProxyNode, 0)

	for _, subscribe := range subscribes {
		// 获取节点
		slog.Info("process subscription", "tag", subscribe.Tag, "url", subscribe.URL)
		subscriptionContent, err := fetcher.FetchSubscription(subscribe.URL, subscribe.UserAgent)
		if err != nil {
			slog.Error("request subscription failed", "error", err)
			continue
		}
		if subscriptionContent == "" {
			slog.Warn("empty subscription content, skip")
			continue
		}

		// 使用不同格式尝试解析
		var subType string
		list, err := clash.Parse(subscriptionContent, subscribe.Tag)
		subType = clash.SubType()

		if err != nil {
			slog.Error("parse subscription failed", "error", err)
			continue
		}
		// 处理节点
		if config.GetConfig().Prefix {
			list = addPrefix(list, subscribe)
		}
		if config.GetConfig().Emoji {
			list = addEmoji(list)
		} else {
			list = removeEmoji(list)
		}
		for i := range list {
			list[i].SubType = subType
		}
		nodes = append(nodes, list...)
	}

	// 处理重复节点名称
	nodes = processDuplicateNodeTag(nodes)

	return nodes, nil
}

func addPrefix(nodes []types.ProxyNode, subscribe config.Subscription) []types.ProxyNode {
	if subscribe.Prefix == "" {
		return nodes
	}

	for i, node := range nodes {
		nodes[i].Tag = subscribe.Prefix + node.Tag
	}
	return nodes
}

func addEmoji(nodes []types.ProxyNode) []types.ProxyNode {
	for i := range nodes {
		nodes[i].Tag = utils.AddNodeEmoji(nodes[i].Tag)
	}
	return nodes
}

func removeEmoji(nodes []types.ProxyNode) []types.ProxyNode {
	for i := range nodes {
		nodes[i].Tag = utils.RemoveNodeEmoji(nodes[i].Tag)
	}
	return nodes
}

func extractTagBase(tag string) string {
	lastDigitPos := len(tag)
	for i := len(tag) - 1; i >= 0; i-- {
		if tag[i] < '0' || tag[i] > '9' {
			break
		}
		lastDigitPos = i
	}
	return tag[:lastDigitPos]
}

func findNextAvailableNumber(base string, existingTags mapset.Set[string]) int {
	num := 1
	for {
		newTag := fmt.Sprintf("%s%d", base, num)
		if !existingTags.Contains(newTag) {
			return num
		}
		num++
	}
}

func processDuplicateNodeTag(nodes []types.ProxyNode) []types.ProxyNode {
	tagSet := mapset.NewSet[string]()

	for i := range nodes {
		originalTag := nodes[i].Tag
		base := extractTagBase(originalTag)

		if !tagSet.Contains(originalTag) {
			tagSet.Add(originalTag)
			continue
		}

		nextNum := findNextAvailableNumber(base, tagSet)
		newTag := fmt.Sprintf("%s%d", base, nextNum)
		nodes[i].Tag = newTag
		tagSet.Add(newTag)
	}

	return nodes
}
