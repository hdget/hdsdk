package utils

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
	"unicode"
)

// RecordErrorStack 将错误信息保存到错误日志文件中
func RecordErrorStack(app string) {
	filename := fmt.Sprintf("%s.err", app)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	defer func() {
		if err == nil {
			file.Close()
		}
	}()
	if err != nil {
		log.Printf("RecordErrorStack: error open file, filename: %s, err: %v", filename, err)
		return
	}

	data := bytes.NewBufferString("=== " + time.Now().String() + " ===\n")
	data.Write(debug.Stack())
	data.WriteString("\n")
	_, err = file.Write(data.Bytes())
	if err != nil {
		log.Printf("RecordErrorStack: error write file, filename: %s, err: %v", filename, err)
	}
}

// ReverseInt64Slice 将[]int64 slice倒序重新排列
func ReverseInt64Slice(numbers []int64) []int64 {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

// GetSliceData 将传过来的数据转换成[]interface{}
func GetSliceData(data interface{}) []interface{} {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil
	}

	sliceLenth := v.Len()
	sliceData := make([]interface{}, sliceLenth)
	for i := 0; i < sliceLenth; i++ {
		sliceData[i] = v.Index(i).Interface()
	}

	return sliceData
}

// GetPagePositions 获取分页的起始值列表
// @return 返回一个二维数组， 第一维是多少页，第二维是每页[]int{start, end}
// e,g: 假设11个数的列表，分页pageSize是5，那么返回的是：
//
//	[]int{
//	   []int{0, 5},
//	   []int{5, 10},
//	   []int{10, 11},
//	}
func GetPagePositions(data interface{}, pageSize int) [][]int {
	listData := GetSliceData(data)
	if listData == nil {
		return nil
	}

	total := len(listData)
	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	pages := make([][]int, 0)
	for i := 0; i < totalPage; i++ {
		start := i * pageSize
		end := (i + 1) * pageSize
		if end > total {
			end = total
		}

		p := []int{start, end}
		pages = append(pages, p)
	}
	return pages
}

// GenerateRandString 生成随机字符串
func GenerateRandString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
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

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// ToFixed 浮点数到指定小数位
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

// GetVarName 获取变量的名字
func GetVarName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

// GetNeo4jPathPattern 解析Neo4j语法的Variable-length pattern
func GetNeo4jPathPattern(args ...int32) string {
	start := int32(-1)
	end := int32(-1)
	switch len(args) {
	case 1:
		start = args[0]
	case 2:
		start = args[0]
		end = args[1]
	}

	expr := "*"
	if start >= 0 {
		if end >= start {
			expr = fmt.Sprintf("*%d..%d", start, end)
		} else {
			expr = fmt.Sprintf("*%d..", start)
		}
	}
	return expr
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

// GetFuncName 从函数实例获取函数名
func GetFuncName(fn any) string {
	tokens := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
	return strings.Split(tokens[len(tokens)-1], "-")[0]
}

func GetStructName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
