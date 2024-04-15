package sqltemplate

import (
	"fmt"
	"github.com/hdget/hdsdk/v2"
	"github.com/hdget/hdsdk/v2/helper/pagination"
	"github.com/hdget/hdsdk/v2/protobuf"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
)

// SqlTemplate 基于模板创建, 只支持where和limit
type SqlTemplate interface {
	With(template string) SqlTemplate
	Limit(listParam *protobuf.ListParam) SqlTemplate
	Next(pkName string, nextParam *protobuf.NextParam) SqlTemplate
	OrderBy(orderBy string) SqlTemplate
	GroupBy(groupBy string) SqlTemplate
	InsertArgs(extraArgs ...any) SqlTemplate
	AppendArgs(extraArgs ...any) SqlTemplate
	JoinSubQuery(tpl SqlTemplate) SqlTemplate

	Where(condition string, args ...any)
	Generate() (string, []any, error)
	Get(v any) error
	Select(v any) error
	Count() (int64, error)
	ExportConditions() ([]string, []any) // 导出where条件
	ImportConditions([]string, []any)    // 导入where条件
}

type baseQuery struct {
	concrete    SqlTemplate
	template    string
	wheres      []string
	args        []any
	limitClause string
	orderBy     string
	groupBy     string
}

type joinSubQuery struct {
	*baseQuery
	alias string
	on    string
}

// SqlTemplate 基于模板创建, 只支持where和limit
type query struct {
	*baseQuery
	joinSubQuerys []*joinSubQuery
}

const (
	defaultNextClause  = "LIMIT 10"
	defaultLimitClause = "LIMIT 0, 10"
)

func NewSqlTemplate() SqlTemplate {
	me := &query{
		baseQuery: &baseQuery{
			wheres: make([]string, 0),
			args:   make([]any, 0),
		},
	}
	me.concrete = me
	return me
}

func NewJoinSubQuery(template, alias, on string) SqlTemplate {
	me := &joinSubQuery{
		baseQuery: &baseQuery{
			template: template,
			wheres:   make([]string, 0),
			args:     make([]any, 0),
		},
		alias: alias,
		on:    on,
	}
	me.concrete = me
	return me
}

/* query */

func (q *query) JoinSubQuery(tpl SqlTemplate) SqlTemplate {
	if jsq, ok := tpl.(*joinSubQuery); ok {
		q.joinSubQuerys = append(q.joinSubQuerys, jsq)
	}
	return q
}

// Generate 最终生成SQL语句
func (q *query) Generate() (string, []any, error) {
	parts := []string{q.template}

	// join子查询
	for _, sq := range q.joinSubQuerys {
		subQuerySql, subQueryArgs, err := sq.Generate()
		if err != nil {
			return "", nil, err
		}

		parts = append(parts, fmt.Sprintf("INNER JOIN (%s) AS %s ON %s", subQuerySql, sq.alias, sq.on))
		q.InsertArgs(subQueryArgs...)
	}

	whereClause := q.getWhereClause()
	if whereClause != "" {
		parts = append(parts, whereClause)
	}

	if q.groupBy != "" {
		parts = append(parts, q.groupBy)
	}

	if q.orderBy != "" {
		parts = append(parts, q.orderBy)
	}

	if q.limitClause != "" {
		parts = append(parts, q.limitClause)
	}

	return q.process(strings.Join(parts, " "))
}

/* joinSubQuery */

func (sq *joinSubQuery) Generate() (string, []any, error) {
	parts := []string{sq.template}

	whereClause := sq.getWhereClause()
	if whereClause != "" {
		parts = append(parts, whereClause)
	}

	if sq.groupBy != "" {
		parts = append(parts, sq.groupBy)
	}

	if sq.orderBy != "" {
		parts = append(parts, sq.orderBy)
	}

	if sq.limitClause != "" {
		parts = append(parts, sq.limitClause)
	}

	return sq.process(strings.Join(parts, " "))
}

/* baseQuery */

func (b *baseQuery) Where(condition string, args ...any) {
	b.wheres = append(b.wheres, fmt.Sprintf("(%s)", condition))
	b.args = append(b.args, args...)
}

func (b *baseQuery) ExportConditions() ([]string, []any) {
	return b.wheres, b.args
}

func (b *baseQuery) ImportConditions(wheres []string, args []any) {
	b.wheres, b.args = wheres, args
}

func (b *baseQuery) With(template string) SqlTemplate {
	b.template = template
	return b.concrete
}

func (b *baseQuery) Limit(listParam *protobuf.ListParam) SqlTemplate {
	if listParam != nil {
		b.limitClause = pagination.New(listParam).GetLimitClause()
		return b.concrete
	}
	b.limitClause = defaultLimitClause
	return b.concrete
}

func (b *baseQuery) Next(pkName string, nextParam *protobuf.NextParam) SqlTemplate {
	if nextParam == nil {
		b.limitClause = defaultNextClause
		return b.concrete
	}

	if nextParam.LastPk > 0 {
		switch nextParam.Direction {
		case protobuf.SortDirection_Desc:
			b.Where(fmt.Sprintf("%s < %d", pkName, nextParam.LastPk))
		default:
			b.Where(fmt.Sprintf("%s > %d", pkName, nextParam.LastPk))
		}
	}

	b.limitClause = fmt.Sprintf("LIMIT %d", nextParam.PageSize)
	return b.concrete
}

func (b *baseQuery) OrderBy(orderBy string) SqlTemplate {
	b.orderBy = fmt.Sprintf("ORDER BY %s", orderBy)
	return b.concrete
}

func (b *baseQuery) GroupBy(groupBy string) SqlTemplate {
	b.groupBy = fmt.Sprintf("GROUP BY %s", groupBy)
	return b.concrete
}

func (b *baseQuery) InsertArgs(extraArgs ...any) SqlTemplate {
	b.args = append(extraArgs, b.args...)
	return b.concrete
}

func (b *baseQuery) AppendArgs(extraArgs ...any) SqlTemplate {
	b.args = append(b.args, extraArgs...)
	return b.concrete
}

func (b *baseQuery) JoinSubQuery(tpl SqlTemplate) SqlTemplate {
	return b.concrete.JoinSubQuery(tpl)
}

func (b *baseQuery) Generate() (string, []any, error) {
	return b.concrete.Generate()
}

func (b *baseQuery) ToSql() (string, []any, error) {
	return b.concrete.Generate()
}

func (b *baseQuery) Get(v any) error {
	xquery, xargs, err := b.concrete.Generate()
	if err != nil {
		return errors.Wrap(err, "generate sql")
	}

	err = hdsdk.Db().My().Get(v, xquery, xargs...)
	if err != nil {
		return errors.Wrap(err, "db get")
	}
	return nil
}

func (b *baseQuery) Select(v any) error {
	xquery, xargs, err := b.concrete.Generate()
	if err != nil {
		return errors.Wrap(err, "generate sql")
	}

	err = hdsdk.Db().My().Select(v, xquery, xargs...)
	if err != nil {
		return errors.Wrap(err, "db select")
	}
	return nil
}

func (b *baseQuery) Count() (int64, error) {
	// 获取total
	xquery, xargs, err := b.concrete.Generate()
	if err != nil {
		return 0, errors.Wrap(err, "generate sql")
	}

	var total int64
	err = hdsdk.Db().My().Get(&total, xquery, xargs...)
	if err != nil {
		return 0, errors.Wrap(err, "db count")
	}
	return total, nil
}

func (b *baseQuery) getWhereClause() string {
	if len(b.wheres) == 0 {
		return ""
	}
	return fmt.Sprintf("WHERE %s", strings.Join(b.wheres, " AND "))
}

// process 后期处理，现在暂时只处理SQL中的IN关键字
func (b *baseQuery) process(query string) (string, []any, error) {
	// 发现如果where里面有IN, 需要特殊处理
	if strings.Contains(strings.ToUpper(query), " IN ") {
		newQuery, newArgs, err := sqlx.In(query, b.args...)
		if err != nil {
			return "", nil, errors.Wrap(err, "sqlx.In")
		}
		newQuery = hdsdk.Db().My().Rebind(newQuery)
		return newQuery, newArgs, nil
	}
	return query, b.args, nil
}
