package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/hdget/hdsdk/lib/pagination"
)

type PaginationParam struct {
	// 指定页面
	Page int64 `form:"page" json:"page"`
	// 指定PageSize
	PageSize int64 `form:"page_size" json:"page_size"`
	// 如果前端传过来了上次页面中的最后的主键id,则要记录起来，用作翻页优化
	LastPk int64 `form:"last_pk" json:"last_pk"`
}

const DefaultPageSize = 10

// GetPagination get *Pagination
// args[0] is default page size
// then you can use pagination.Paging
// or pagination.GetSQLClause(total) 去获取分页的limit子句
func GetPagination(c *gin.Context) *pagination.Pagination {
	var p PaginationParam
	if c.ShouldBind(&p) != nil {
		return nil
	}

	return pagination.New(p.Page, p.PageSize)
}
