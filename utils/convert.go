package utils

import (
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"strings"
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

// String 尝试将值转换成字符串
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

// CsvToInt64s 将逗号分隔的string尝试转换成[1,2,3...]的int64 slice
// Csv means Comma Separated Value
func CsvToInt64s(strValue string) []int64 {
	if len(strValue) == 0 {
		return nil
	}

	tokens := strings.Split(strValue, ",")
	if len(tokens) == 0 {
		return nil
	}

	return ToInt64Slice(tokens)
}

// Int64sToCsv 将int64 slice转换成用逗号分隔的字符串: 1,2,3
func Int64sToCsv(int64s []int64) string {
	return strings.Join(ToStringSlice(int64s), ",")
}

// Int32sToCsv 将int32 slice转换成用逗号分隔的字符串: 1,2,3
func Int32sToCsv(int32s []int64) string {
	return strings.Join(ToStringSlice(int32s), ",")
}

// CsvToInt32s 将逗号分隔的string尝试转换成[1,2,3...]的int64 slice
// Csv means Comma Separated Value
func CsvToInt32s(strValue string) []int64 {
	if len(strValue) == 0 {
		return nil
	}

	tokens := strings.Split(strValue, ",")
	if len(tokens) == 0 {
		return nil
	}

	return ToInt32Slice(tokens)
}

// ToStringSlice 将int64 slice转换成["1", "2", "3"...]字符串slice
func ToStringSlice(int64Slice []int64) []string {
	if len(int64Slice) == 0 {
		return nil
	}

	stringList := make([]string, 0)
	for _, item := range int64Slice {
		stringList = append(stringList, strconv.FormatInt(item, 10))
	}

	return stringList
}

// ToInt64Slice 将string slice转换成[1,2,3...]的int64 slice
func ToInt64Slice(strSlice []string) []int64 {
	if len(strSlice) == 0 {
		return nil
	}
	stringList := make([]int64, 0)
	for _, item := range strSlice {
		parseInt, _ := strconv.ParseInt(item, 10, 64)
		stringList = append(stringList, parseInt)
	}

	return stringList
}

// ToInt32Slice 将string slice转换成[1,2,3...]的int32 slice
func ToInt32Slice(strSlice []string) []int64 {
	if len(strSlice) == 0 {
		return nil
	}
	stringList := make([]int64, 0)
	for _, item := range strSlice {
		parseInt, _ := strconv.ParseInt(item, 10, 32)
		stringList = append(stringList, parseInt)
	}

	return stringList
}
