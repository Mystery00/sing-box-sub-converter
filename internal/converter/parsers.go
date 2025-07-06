package converter

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"log/slog"
	"sing-box-sub-converter/internal/config"
	"sing-box-sub-converter/internal/converter/types"
	"sing-box-sub-converter/internal/fetcher"
	"sing-box-sub-converter/utils"
	"strings"
	"time"
)

// ProcessSubscribes 处理所有订阅
func ProcessSubscribes(subscribes []config.Subscription) ([]types.ProxyNode, error) {
	nodes := make([]types.ProxyNode, 0)

	for _, subscribe := range subscribes {
		// 获取节点
		slog.Info("处理订阅", "tag", subscribe.Tag, "url", subscribe.URL)
		subscriptionContent, subInfo, err := fetcher.FetchSubscription(subscribe.URL, subscribe.UserAgent)
		if err != nil {
			slog.Error("请求订阅失败", "error", err)
			continue
		}
		if subscriptionContent == "" {
			slog.Warn("订阅内容为空，跳过")
			continue
		}

		// 使用不同格式尝试解析
		var subType string
		var list []types.ProxyNode

		for _, parser := range parsers() {
			list, err = parser.Parse(subscriptionContent, subscribe.Tag)
			subType = parser.SubType()
			if err == nil {
				break
			}
			slog.Info(fmt.Sprintf("使用%s解析订阅失败，尝试其他方式", subType))
		}

		if err != nil {
			slog.Error("解析订阅失败", "error", err)
			continue
		}
		// 处理节点
		if config.GetConfig().Emoji {
			list = addEmoji(list)
		} else {
			list = removeEmoji(list)
		}
		if config.GetConfig().Prefix {
			list = addPrefix(list, subscribe)
		}
		for i := range list {
			list[i].SubType = subType
		}
		nodes = append(nodes, list...)
		if config.GetConfig().ShowSubInNodes && subInfo != nil {
			remainBytes := bytesToGB(subInfo.Total - subInfo.Upload - subInfo.Download)
			remainDays := calculateRemainDays(subInfo.Expire)
			prefix := strings.TrimSpace(subscribe.Prefix)
			if prefix == "" {
				prefix = strings.TrimSpace(subscribe.Tag)
			}
			nodes = append(nodes, types.ProxyNode{
				Type:        types.ProxyNodeTypeSubInfo,
				Tag:         fmt.Sprintf("%s 剩余流量：%s 剩余天数：%s", prefix, remainBytes, remainDays),
				Address:     "subInfo",
				Port:        "1",
				FromSub:     subscribe.Tag,
				SubType:     "",
				ProxyDetail: nil,
			})
		}
	}

	// 处理重复节点名称
	nodes = processDuplicateNodeTag(nodes)

	return nodes, nil
}

func calculateRemainDays(timeInMills int64) string {
	expireTime := time.Unix(timeInMills, 0)
	remainDays := time.Until(expireTime).Hours() / 24
	return fmt.Sprintf("%d天", int(remainDays))
}

func bytesToGB(bytes int64) string {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	return fmt.Sprintf("%.2f GB", gb)
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
