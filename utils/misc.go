package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"
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

// GetRealIP 获取真实IP
func GetRealIP(c *gin.Context) string {
	xForwardInfo := c.GetHeader("X-Forwarded-For")
	if xForwardInfo != "" {
		ips := strings.Split(xForwardInfo, ",")
		// 返回第一个真实IP
		if len(ips) >= 1 {
			return ips[0]
		}
	}
	return c.ClientIP()
}
