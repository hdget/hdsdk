package pagination

import (
	"fmt"
	"math"
	"reflect"
)

const DefaultPageSize = 10

type Pagination struct {
	page     int64 // 第几页
	pageSize int64 // 每页几项
}

// New pagination
// args[0] is default page size
func New(page, pageSize int64, args ...int64) *Pagination {
	// 处理当前页面
	if page == 0 {
		page = 1
	}

	// 处理每页大小
	if pageSize <= 0 {
		pageSize = DefaultPageSize
		if len(args) > 0 { // if specified default page size, then use it
			pageSize = args[0]
		}
	}

	return &Pagination{page: page, pageSize: pageSize}
}

// Paging 分页
// @return total
// @return []interface{} 分页后的数据
func (p *Pagination) Paging(data interface{}) (int64, []interface{}) {
	sliceData := convertToSlice(data)
	total := int64(len(sliceData))
	start, end := GetStartEndPosition(p.page, p.pageSize, total)
	return total, sliceData[start:end]
}

// GetLimitClause 获取limit sql子句
func (p *Pagination) GetLimitClause() string {
	if p == nil || p.page == 0 {
		return fmt.Sprintf("LIMIT %d", p.pageSize)
	}

	start := (p.page - 1) * p.pageSize
	return fmt.Sprintf("LIMIT %d, %d", start, p.pageSize)
}

// GetSQLClause 获取翻页SQL查询语句
//
// 1. 假如前端没有传过来last_pk, 那么返回值是 last_pk, LIMIT子句(LIMIT offset, pageSize)
// e,g: 0, "LIMIT 20, 10" => 在数据库查询时可能会被组装成 WHERE pk > 0 ...  LIMIT 20, 10
//
// 2. 假如前端传过来了last_pk, 那么返回值是 last_pk, LIMIT子句(LIMIT pageSize)
// e,g: 123,"LIMIT 10" => 在数据库查询时可能会被组装成 WHERE pk > 123 ...  LIMIT 10
//func (p *Pagination) GetSQLClause(total int64) string {
//	if p == nil {
//		return ""
//	}
//
//	// 如果total值为0, 默认返回指定页面
//	if total == 0 {
//		return "LIMIT 0"
//	}
//
//	start := (p.page - 1) * p.pageSize
//	return fmt.Sprintf("LIMIT %d, %d", start, p.pageSize)
//	//start, end := GetStartEndPosition(p.page, p.pageSize, total)
//	//
//	//return fmt.Sprintf("LIMIT %d, %d", start, end-start)
//}

// CalculatePages 计算页面，获取带有起始值的页面的数组
// @return 返回一个二维数组， 第一维是多少页，第二维是每页[]int64{start, end}
// e,g: 假设11个数的列表，分页pageSize是5，那么返回的是：
// []int64{
//    []int64{0, 5},
//    []int64{5, 10},
//    []int64{10, 11},
// }
func CalculatePages(total, pageSize int64) [][]int64 {
	totalPage := int64(math.Ceil(float64(total) / float64(pageSize)))

	pages := make([][]int64, 0)
	for i := int64(0); i < totalPage; i++ {
		start := i * pageSize
		end := (i + 1) * pageSize
		if end > total {
			end = total
		}

		p := []int64{start, end}
		pages = append(pages, p)
	}
	return pages
}

// GetStartEndPosition 如果是按列表slice进行翻页的话， 计算slice的起始index
func GetStartEndPosition(page, pageSize, total int64) (int64, int64) {
	start := (page - 1) * pageSize
	end := page * pageSize

	if end > total {
		end = total
	}

	if start > end {
		start = end
	}

	return start, end
}

// convertToSlice convert interface{} to []interface{}
func convertToSlice(data interface{}) []interface{} {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil
	}

	sliceLength := v.Len()
	sliceData := make([]interface{}, sliceLength)
	for i := 0; i < sliceLength; i++ {
		sliceData[i] = v.Index(i).Interface()
	}

	return sliceData
}
