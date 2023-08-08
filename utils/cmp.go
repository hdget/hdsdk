package utils

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
