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
// []int{
//    []int{0, 5},
//    []int{5, 10},
//    []int{10, 11},
// }
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
