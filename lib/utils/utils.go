package utils

import (
	"github.com/hdget/hdsdk/v1/protobuf"
	"github.com/hdget/hdutils"
)

const (
	DefaultPageSize = 10
)

var (
	defaultPageParam = &protobuf.ListParam{
		Page:     1,
		PageSize: DefaultPageSize,
	}
)

func NewWithParam(listParam *protobuf.ListParam) hdutils.Pagination {
	p := listParam
	if p == nil {
		p = defaultPageParam
	}

	// 处理当前页面
	page := p.Page
	if page == 0 {
		page = 1
	}

	// 处理每页大小
	pageSize := p.PageSize
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	return hdutils.NewPagination(page, pageSize)
}
