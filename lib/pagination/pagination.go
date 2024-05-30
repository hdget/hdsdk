package pagination

import (
	"github.com/hdget/hdsdk/v2/protobuf"
	"github.com/hdget/hdutils/page"
)

const (
	DefaultPageSize = 10
)

var (
	defaultPageParam = &protobuf.ListParam{
		Page:     1,
		PageSize: page.DefaultPageSize,
	}
)

func New(listParam *protobuf.ListParam) page.Pagination {
	p := listParam
	if p == nil {
		p = defaultPageParam
	}

	// 处理当前页面
	currentPage := p.Page
	if currentPage == 0 {
		currentPage = 1
	}

	// 处理每页大小
	pageSize := p.PageSize
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	return page.NewPagination(currentPage, pageSize)
}
