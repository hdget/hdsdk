package utils

import (
	"strings"
	"unicode"
)

// CleanString 处理字符串, args[0]为是否转换为小写
func CleanString(origStr string, args ...bool) string {
	// 1. 去除前后空格
	s := strings.TrimSpace(origStr)

	// 2. 是否转换小写
	toLower := false
	if len(args) > 0 {
		toLower = args[0]
	}

	if toLower {
		s = strings.ToLower(s)
	}

	// 去除不可见字符
	s = RemoveInvisibleCharacter(s)
	return s
}

// RemoveInvisibleCharacter 去除掉不能显示的字符
func RemoveInvisibleCharacter(origStr string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, origStr)
}
