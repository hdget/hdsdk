package db

import (
	"github.com/Masterminds/squirrel"
	"github.com/hdget/hdsdk/lib/pagination"
	"github.com/hdget/hdsdk/protobuf"
	"github.com/hdget/hdsdk/types"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type BaseDbProvider struct {
	Default *BaseDbClient            // 缺省的数据库连接
	Main    *BaseDbClient            // 主数据库连接,可读写
	Slaves  []*BaseDbClient          // 只读从数据库列表，只可读
	Items   map[string]*BaseDbClient // 额外数据库
}

func (b *BaseDbProvider) My() types.DbClient {
	return b.Default
}

func (b *BaseDbProvider) Master() types.DbClient {
	return b.Main
}

func (b *BaseDbProvider) Slave(i int) types.DbClient {
	return b.Slaves[i]
}

func (b *BaseDbProvider) By(s string) types.DbClient {
	return b.Items[s]
}

type BaseDbClient struct {
	Db       *sqlx.DB // 缺省的数据库连接
	_builder squirrel.Sqlizer
}

func (b *BaseDbClient) Rebind(query string) string {
	return b.Db.Rebind(query)
}

func (b *BaseDbClient) HGet(v any) error {
	xquery, xargs, err := b._builder.ToSql()
	if err != nil {
		return err
	}
	return b.Db.Get(v, xquery, xargs...)
}

func (b *BaseDbClient) HSelect(v any, args ...*protobuf.ListParam) error {
	var p *pagination.Pagination
	bd := b._builder
	if len(args) > 0 {
		selBd, ok := bd.(squirrel.SelectBuilder)
		if !ok {
			return errors.New("invalid select builder")
		}

		p = pagination.NewWithParam(args[0])
		bd = selBd.Offset(p.Offset).Limit(p.PageSize)
	}

	xquery, xargs, err := bd.ToSql()
	if err != nil {
		return err
	}
	return b.Db.Select(v, xquery, xargs...)
}

func (b *BaseDbClient) HCount() (int64, error) {
	xquery, xargs, err := b._builder.ToSql()
	if err != nil {
		return 0, err
	}

	var total int64
	err = b.Db.Get(&total, xquery, xargs...)
	return total, err
}

func (b *BaseDbClient) HQuery(args ...*protobuf.ListParam) (*sqlx.Rows, error) {
	var p *pagination.Pagination
	bd := b._builder
	if len(args) > 0 {
		selBd, ok := bd.(squirrel.SelectBuilder)
		if !ok {
			return nil, errors.New("invalid select builder")
		}

		p = pagination.NewWithParam(args[0])
		bd = selBd.Offset(p.Offset).Limit(p.PageSize)
	}

	xquery, xargs, err := bd.ToSql()
	if err != nil {
		return nil, err
	}
	return b.Db.Queryx(xquery, xargs...)
}

func (b *BaseDbClient) Select(dest any, query string, args ...any) error {
	return b.Db.Select(dest, query, args...)
}

func (b *BaseDbClient) Get(dest any, query string, args ...any) error {
	return b.Db.Get(dest, query, args...)
}

func (b *BaseDbClient) Queryx(query string, args ...any) (*sqlx.Rows, error) {
	return b.Db.Queryx(query, args...)
}

func (b *BaseDbClient) ToSql() (string, []any, error) {
	return b._builder.ToSql()
}

func (b *BaseDbClient) Sq(builder squirrel.SelectBuilder) types.DbClient {
	b._builder = builder
	return b
}
