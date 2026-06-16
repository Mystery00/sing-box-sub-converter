package utils

import "maps"

// DeepMerge 以 base 为基底，用 override 覆盖，返回合并后的新 map。
// override 字段优先；同名嵌套 map 递归合并；数组及标量整体覆盖。
// 不修改入参。
func DeepMerge(base, override map[string]any) map[string]any {
	result := make(map[string]any, len(base))
	maps.Copy(result, base)
	for k, ov := range override {
		if bv, exist := result[k]; exist {
			bm, bok := bv.(map[string]any)
			om, ook := ov.(map[string]any)
			if bok && ook {
				result[k] = DeepMerge(bm, om)
				continue
			}
		}
		result[k] = ov
	}
	return result
}
