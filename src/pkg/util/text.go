package util

import (
	"fmt"
	"strings"
)

// SanitizeRawText 清理文本中不支持的字符
// 保留ASCII可打印字符(32-126)和基本控制字符(\r\n\t)
// 其他字符将被替换为其Unicode编码表示 [0xXXXX]
func SanitizeRawText(text string) string {
	if text == "" {
		return ""
	}

	// 创建一个新的字符串builder
	var result strings.Builder
	result.Grow(len(text))

	// 遍历字符串中的每个字符
	for _, r := range text {
		// 保留ASCII可打印字符(32-126)和基本控制字符(\r\n\t)
		if (r >= 32 && r <= 126) || r == '\r' || r == '\n' || r == '\t' {
			result.WriteRune(r)
		} else {
			// 对于不支持的字符，替换为其Unicode编码表示
			result.WriteString(fmt.Sprintf("[0x%X]", r))
		}
	}

	return result.String()
}
