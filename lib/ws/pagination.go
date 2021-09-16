package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

type PaginationParam struct {
	// 指定页面
	Page int `form:"page" json:"page"`
	// 指定PageSize
	PageSize int `form:"page_size" json:"page_size"`
	// 如果前端传过来了上次页面中的最后的主键id,则要记录起来，用作翻页优化
	LastPk int64 `form:"last_pk" json:"last_pk"`
}

const DefaultPageSize = 10

// GetPaginationParam obtain pagination params
// if specified args, then args[0] is specified default page size
func GetPaginationParam(c *gin.Context, args ...int) *PaginationParam {
	var p PaginationParam
	if c.ShouldBind(&p) != nil {
		return nil
	}

	if p.Page == 0 {
		p.Page = 1
	}

	if p.PageSize <= 0 {
		p.PageSize = DefaultPageSize
		if len(args) > 0 { // if specified default page size, then use it
			p.PageSize = args[0]
		}
	}

	return &p
}

// Paging 将数据分页
// @return total
// @return []interface{}
func Paging(c *gin.Context, data interface{}) (int64, []interface{}) {
	sliceData := convertToSlice(data)
	if len(sliceData) == 0 {
		return 0, nil
	}

	// 分页数据
	param := GetPaginationParam(c)
	total := int64(len(sliceData))
	// slice
	start, end := getStartEndPosition(param, total)
	return total, sliceData[start:end]
}

// GetPagingSQLClause 获取翻页SQL查询语句
//
// 1. 假如前端没有传过来last_pk, 那么返回值是 last_pk, LIMIT子句(LIMIT offset, pageSize)
// e,g: 0, "LIMIT 20, 10" => 在数据库查询时可能会被组装成 WHERE pk > 0 ...  LIMIT 20, 10
//
// 2. 假如前端传过来了last_pk, 那么返回值是 last_pk, LIMIT子句(LIMIT pageSize)
// e,g: 123,"LIMIT 10" => 在数据库查询时可能会被组装成 WHERE pk > 123 ...  LIMIT 10

func GetPagingSQLClause(p *PaginationParam, total int64) (int64, string) {
	if p == nil {
		return 0, fmt.Sprintf("LIMIT 0, %d", p.PageSize)
	}

	// 如果前端没有传过来total值，则GetPageSQLInfo必须在args中传递total值，否则报错
	if total == 0 {
		return 0, fmt.Sprintf("LIMIT 0, %d", p.PageSize)
	}

	if p.LastPk == 0 {
		start, end := getStartEndPosition(p, total)
		limitClause := fmt.Sprintf("LIMIT %d, %d", start, end-start)
		return 0, limitClause
	}

	// 如果传过来了lastPk, 则只构成LIMIT SIZE子句
	return p.LastPk, fmt.Sprintf("LIMIT %d", p.PageSize)
}

// getStartEndPosition 如果是按列表slice进行翻页的话， 计算slice的起始index
// args[0] 前端传过来的total值
func getStartEndPosition(p *PaginationParam, total int64) (int64, int64) {
	if p == nil {
		return 0, 0
	}

	start := int64((p.Page - 1) * p.PageSize)
	end := int64(p.Page * p.PageSize)

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
