package strkit

import "strings"

// 传入字符串、分隔符，取第 index 段，index 从 1 开始
func GetSegment(str, sep string, index int) string {
	// 使用指定的分隔符分割字符串
	parts := strings.Split(str, sep)

	// 将 index 转换为从 1 开始的逻辑，检查是否在有效范围内
	if index > 0 && index <= len(parts) {
		return parts[index-1]
	}

	// 如果 index 无效，返回空字符串
	return ""
}

// 传入字符串、分隔符，取最后一段
func GetLastSegment(str, sep string) string {
	parts := strings.Split(str, sep)
	lastPart := parts[len(parts)-1]
	return lastPart
}
