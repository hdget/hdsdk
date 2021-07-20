package utils

import (
	jsoniter "github.com/json-iterator/go"
	"unsafe"
)

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// 尝试将值转换成字符串
func String(value interface{}) (string, error) {
	switch reply := value.(type) {
	case string:
		return reply, nil
	}

	bs, err := jsoniter.Marshal(value)
	if err != nil {
		return "", err
	}
	return BytesToString(bs), nil
}
