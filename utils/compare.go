package utils

// IntSliceContains 检查整型slice中是否含有
func IntSliceContains(list []int, checkItem int) bool {
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

// Int64SliceContains 检查整型slice中是否含有
func Int64SliceContains(list []int64, checkItem int64) bool {
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

// StringSliceContains 检查字符串slice中是否含有
func StringSliceContains(list []string, checkItem string) bool {
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

// IsAlphanumeric check if the string contains only letters and numbers. Empty string is valid.
func IsAlphanumeric(s string) bool {
	for _, v := range s {
		if ('Z' < v || v < 'A') && ('z' < v || v < 'a') && ('9' < v || v < '0') {
			return false
		}
	}
	return true
}
