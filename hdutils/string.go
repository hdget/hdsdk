package hdutils

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	nonNumericRegex      = regexp.MustCompile(`[^0-9 ]+`)       // 非数字
	nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`) // 非英文字符和数字
	nonChineseRegex      = regexp.MustCompile(`[^\p{Han}]+`)    // 非汉字
)

// Numeric 去除数字以外的所有字符
func Numeric(s string) string {
	return nonNumericRegex.ReplaceAllString(s, "")
}

// AlphaNumeric 去除字符数字以外的所有字符
func AlphaNumeric(s string) string {
	return nonAlphanumericRegex.ReplaceAllString(s, "")
}

// Chinese 去除中文以外的所有字符
func Chinese(s string) string {
	return nonChineseRegex.ReplaceAllString(s, "")
}

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
