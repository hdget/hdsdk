package sql

import (
	"fmt"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/lib/pagination"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
)

// SqlBuilder 查询SQL创建支持所有语法
type SqlBuilder struct {
	wheres        []string
	args          []any
	hasInSymbol   bool // 是否条件里面有"in"
	limitClause   string
	selectClause  string
	fromClause    string
	joins         []string
	groupByClause string
	orderByClause string
}

func NewSqlBuilder() *SqlBuilder {
	return &SqlBuilder{
		wheres: make([]string, 0),
		args:   make([]any, 0),
		joins:  make([]string, 0),
	}
}

func (b *SqlBuilder) Count(args ...string) *SqlBuilder {
	keyword := "1"
	if len(args) > 0 {
		keyword = args[0]
	}
	b.selectClause = fmt.Sprintf("SELECT COUNT(%s)", keyword)
	return b
}

func (b *SqlBuilder) Select(items ...string) *SqlBuilder {
	b.selectClause = fmt.Sprintf("SELECT %s", strings.Join(items, ","))
	return b
}

func (b *SqlBuilder) From(table string, alias ...string) *SqlBuilder {
	from := fmt.Sprintf("FROM %s", table)
	if len(alias) > 0 {
		from = fmt.Sprintf("%s AS %s", from, alias[0])
	}
	b.fromClause = from
	return b
}

func (b *SqlBuilder) InnerJoin(table, alias, on string, args ...any) *SqlBuilder {
	b.joins = append(b.joins, fmt.Sprintf("INNER JOIN %s AS %s ON %s", table, alias, on))
	b.args = append(b.args, args...)
	return b
}

func (b *SqlBuilder) LeftJoin(table, alias, on string, args ...any) *SqlBuilder {
	b.joins = append(b.joins, fmt.Sprintf("LEFT JOIN %s AS %s ON %s", table, alias, on))
	b.args = append(b.args, args...)
	return b
}

func (b *SqlBuilder) RightJoin(table, alias, on string, args ...any) *SqlBuilder {
	b.joins = append(b.joins, fmt.Sprintf("RIGTH JOIN %s AS %s ON %s", table, alias, on))
	b.args = append(b.args, args...)
	return b
}

func (b *SqlBuilder) GroupBy(groupBy string) *SqlBuilder {
	b.groupByClause = fmt.Sprintf("GROUP BY %s", groupBy)
	return b
}

func (b *SqlBuilder) OrderBy(orderBy string) *SqlBuilder {
	b.orderByClause = fmt.Sprintf("ORDER BY %s", orderBy)
	return b
}

func (b *SqlBuilder) Limit(listParam any) *SqlBuilder {
	if listParam != nil {
		param, ok := listParam.(*pagination.ListParam)
		if ok {
			b.limitClause = pagination.New(param.Page, param.PageSize).GetLimitClause()
			return b
		}
	}
	b.limitClause = defaultLimitClause
	return b
}

func (b *SqlBuilder) Build() (string, []any, error) {
	// 第一个为select,from
	clauses := []string{
		b.selectClause,
		b.fromClause,
	}

	// join
	if len(b.joins) > 0 {
		clauses = append(clauses, strings.Join(b.joins, ""))
	}

	// where
	clauses = append(clauses, b.getWhereClause())

	// IMPORTANT: 顺序很重要
	if b.groupByClause != "" {
		clauses = append(clauses, b.groupByClause)
	}

	if b.orderByClause != "" {
		clauses = append(clauses, b.groupByClause)
	}

	if b.limitClause != "" {
		clauses = append(clauses, b.limitClause)
	}

	return b.process(strings.Join(clauses, " "))
}

func (b *SqlBuilder) getWhereClause() string {
	if len(b.wheres) == 0 {
		return "1=1"
	}
	whereClause := strings.Join(b.wheres, " AND ")
	if strings.Contains(strings.ToUpper(whereClause), " IN ") {
		b.hasInSymbol = true
	}
	return whereClause
}

// process 后期处理，现在暂时只处理SQL中的IN关键字
func (b *SqlBuilder) process(query string) (string, []any, error) {
	// 发现如果where里面有IN, 需要特殊处理
	if b.hasInSymbol {
		newQuery, newArgs, err := sqlx.In(query, b.args...)
		if err != nil {
			return "", nil, errors.Wrap(err, "sqlx.In")
		}
		newQuery = hdsdk.Mysql.My().Rebind(newQuery)
		return newQuery, newArgs, nil
	}
	return query, b.args, nil
}
