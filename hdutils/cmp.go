package hdutils

import (
	"regexp"
	"strings"
)

var (
	regexIsMobile = regexp.MustCompile(`^1{3-9}\d{9}$`)
)

// IsValidMobile check if the string is valid chinese mobile number
func IsValidMobile(mobile string) bool {
	return regexIsMobile.MatchString(mobile)
}

// IsAlphanumeric check if the string contains only letters and numbers. Empty string is valid.
func IsAlphanumeric(s string) bool {
	for _, v := range s {
		if ('Z' < v || v < 'A') && ('z' < v || v < 'a') && ('9' < v || v < '0') {
			return false
		}
	}
	return true
}

// IsNumeric check if the string contains only numbers. Empty string is valid.
func IsNumeric(s string) bool {
	for _, v := range s {
		if '9' < v || v < '0' {
			return false
		}
	}
	return true
}

// IsImageData 是否是图像数据
func IsImageData(data []byte) bool {
	// image formats and magic numbers
	var magicTable = map[string]string{
		"\xff\xd8\xff":      "image/jpeg",
		"\x89PNG\r\n\x1a\n": "image/png",
		"GIF87a":            "image/gif",
		"GIF89a":            "image/gif",
	}
	s := BytesToString(data)
	for magic := range magicTable {
		if strings.HasPrefix(s, magic) {
			return true
		}
	}
	return false
}

func Contains[T comparable](list []T, checkItem T) bool {
	if len(list) == 0 {
		return false
	}

	for _, item := range list {
		if item == checkItem {
			return true
		}
	}

	return false
}
