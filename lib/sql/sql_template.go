package sql

import (
	"fmt"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/lib/pagination"
	"github.com/hdget/hdsdk/protobuf"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
)

// sqlTemplate 基于模板创建, 只支持where和limit
type sqlTemplate struct {
	template    string
	wheres      []string
	args        []any
	hasInSymbol bool // 是否条件里面有"in"
	limitClause string
	orderBy     string
	groupBy     string
}

const (
	defaultNextClause  = "LIMIT 10"
	defaultLimitClause = "LIMIT 0, 10"
)

func NewSqlTemplate() *sqlTemplate {
	return &sqlTemplate{
		wheres: make([]string, 0),
		args:   make([]any, 0),
	}
}

func (h *sqlTemplate) Where(condition string, args ...any) {
	h.wheres = append(h.wheres, fmt.Sprintf("(%s)", condition))
	h.args = append(h.args, args...)
}

func (h *sqlTemplate) With(template string) *sqlTemplate {
	h.template = template
	return h
}

func (h *sqlTemplate) Limit(listParam *protobuf.ListParam) *sqlTemplate {
	if listParam != nil {
		h.limitClause = pagination.New(listParam.Page, listParam.PageSize).GetLimitClause()
		return h
	}
	h.limitClause = defaultLimitClause
	return h
}

func (h *sqlTemplate) Next(pkName string, nextParam *protobuf.NextParam) *sqlTemplate {
	if nextParam == nil {
		h.limitClause = defaultNextClause
		return h
	}

	if nextParam.LastPk > 0 {
		switch nextParam.Direction {
		case protobuf.SortDirection_Desc:
			h.Where(fmt.Sprintf("%s < %d", pkName, nextParam.LastPk))
		default:
			h.Where(fmt.Sprintf("%s > %d", pkName, nextParam.LastPk))
		}
	}

	h.limitClause = fmt.Sprintf("LIMIT %d", nextParam.PageSize)
	return h
}

func (h *sqlTemplate) OrderBy(orderBy string) *sqlTemplate {
	h.orderBy = fmt.Sprintf("ORDER BY %s", orderBy)
	return h
}

func (h *sqlTemplate) GroupBy(groupBy string) *sqlTemplate {
	h.groupBy = fmt.Sprintf("GROUP BY %s", groupBy)
	return h
}

func (h *sqlTemplate) Get(v any) error {
	query, args, err := h.Generate()
	if err != nil {
		return errors.Wrap(err, "build query")
	}

	err = hdsdk.Mysql.My().Get(&v, query, args...)
	if err != nil {
		return errors.Wrap(err, "db get")
	}
	return nil
}

// Generate 最终生成SQL语句
func (h *sqlTemplate) Generate() (string, []any, error) {
	parts := []string{h.template}

	whereClause := h.getWhereClause()
	if whereClause != "" {
		parts = append(parts, whereClause)
	}

	if h.groupBy != "" {
		parts = append(parts, h.groupBy)
	}

	if h.orderBy != "" {
		parts = append(parts, h.orderBy)
	}

	if h.limitClause != "" {
		parts = append(parts, h.limitClause)
	}

	return h.process(strings.Join(parts, " "))
}

func (h *sqlTemplate) Count() (int64, error) {
	// 获取total
	query, args, err := h.Generate()
	if err != nil {
		return 0, errors.Wrap(err, "build query")
	}
	var total int64
	err = hdsdk.Mysql.My().Get(&total, query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db count")
	}
	return total, nil
}

func (h *sqlTemplate) InsertArgs(extraArgs ...any) *sqlTemplate {
	h.args = append(extraArgs, h.args...)
	return h
}

func (h *sqlTemplate) AppendArgs(extraArgs ...any) *sqlTemplate {
	h.args = append(h.args, extraArgs...)
	return h
}

func (h *sqlTemplate) getWhereClause() string {
	if len(h.wheres) == 0 {
		return ""
	}

	whereClause := fmt.Sprintf("WHERE %s", strings.Join(h.wheres, " AND "))
	if strings.Contains(strings.ToUpper(whereClause), " IN ") {
		h.hasInSymbol = true
	}
	return whereClause
}

// process 后期处理，现在暂时只处理SQL中的IN关键字
func (h *sqlTemplate) process(query string) (string, []any, error) {
	// 发现如果where里面有IN, 需要特殊处理
	if h.hasInSymbol {
		newQuery, newArgs, err := sqlx.In(query, h.args...)
		if err != nil {
			return "", nil, errors.Wrap(err, "sqlx.In")
		}
		newQuery = hdsdk.Mysql.My().Rebind(newQuery)
		return newQuery, newArgs, nil
	}
	return query, h.args, nil
}
