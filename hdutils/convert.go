package hdutils

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strings"
	"unicode"
	"unsafe"
)

var (
	rxCameling = regexp.MustCompile(`[\p{L}\p{N}]+`)
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

// ToString 尝试将值转换成字符串
func ToString(value interface{}) (string, error) {
	switch reply := value.(type) {
	case string:
		return reply, nil
	case []byte:
		return BytesToString(reply), nil
	}

	bs, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return BytesToString(bs), nil
}

func ToBytes(value interface{}) ([]byte, error) {
	var data []byte
	switch t := value.(type) {
	case string:
		data = StringToBytes(t)
	case []byte:
		data = t
	default:
		v, err := json.Marshal(value)
		if err != nil {
			return nil, errors.Wrapf(err, "marshal value, value: %v", value)
		}
		data = v
	}
	return data, nil
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

// CsvToInt32s 将逗号分隔的string尝试转换成[1,2,3...]的int32 slice
// Csv means Comma Separated Value
func CsvToInt32s(strValue string) []int32 {
	if len(strValue) == 0 {
		return nil
	}

	tokens := strings.Split(strValue, ",")
	if len(tokens) == 0 {
		return nil
	}

	return ToInt32Slice(tokens)
}

// Int64sToCsv 将int64 slice转换成用逗号分隔的字符串: 1,2,3
func Int64sToCsv(int64s []int64) string {
	return strings.Join(cast.ToStringSlice(int64s), ",")
}

// Int32sToCsv 将int32 slice转换成用逗号分隔的字符串: 1,2,3
func Int32sToCsv(int32s []int32) string {
	return strings.Join(cast.ToStringSlice(int32s), ",")
}

// ToInt64Slice 将string slice转换成[1,2,3...]的int64 slice
func ToInt64Slice(strSlice []string) []int64 {
	if len(strSlice) == 0 {
		return nil
	}
	int64s := make([]int64, 0)
	for _, item := range strSlice {
		int64s = append(int64s, cast.ToInt64(item))
	}
	return int64s
}

// ToInt32Slice 将string slice转换成[1,2,3...]的int32 slice
func ToInt32Slice(strSlice []string) []int32 {
	if len(strSlice) == 0 {
		return nil
	}
	int32s := make([]int32, 0)
	for _, item := range strSlice {
		int32s = append(int32s, cast.ToInt32(item))
	}
	return int32s
}

// ToCamelCase converts from underscore separated form to camel case form.
func ToCamelCase(str string) string {
	byteSrc := []byte(str)
	chunks := rxCameling.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		chunks[idx] = cases.Title(language.English).Bytes(val)
	}
	return string(bytes.Join(chunks, nil))
}

// ToSnakeCase converts from camel case form to underscore separated form.
func ToSnakeCase(s string) string {
	s = ToCamelCase(s)
	runes := []rune(s)
	length := len(runes)
	var out []rune
	for i := 0; i < length; i++ {
		out = append(out, unicode.ToLower(runes[i]))
		if i+1 < length && (unicode.IsUpper(runes[i+1]) && unicode.IsLower(runes[i])) {
			out = append(out, '_')
		}
	}

	return string(out)
}
