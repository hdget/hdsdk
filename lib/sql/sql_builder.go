package sql

import (
	"github.com/Masterminds/squirrel"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/lib/pagination"
	"github.com/hdget/hdsdk/protobuf"
	"github.com/jmoiron/sqlx"
)

type SqlBuilder struct {
	selectBuilder squirrel.SelectBuilder
	_db           *sqlx.DB
}

func NewSqlBuilder(selectBuilder squirrel.SelectBuilder) *SqlBuilder {
	return &SqlBuilder{selectBuilder: selectBuilder}
}

func (b *SqlBuilder) Master() *SqlBuilder {
	b._db = hdsdk.Mysql.Master()
	return b
}

func (b *SqlBuilder) Slave(index int) *SqlBuilder {
	b._db = hdsdk.Mysql.Slave(index)
	return b
}

func (b *SqlBuilder) By(name string) *SqlBuilder {
	b._db = hdsdk.Mysql.By(name)
	return b
}

func (b *SqlBuilder) Get(v any) error {
	xquery, xargs, err := b.selectBuilder.ToSql()
	if err != nil {
		return err
	}
	if b._db != nil {
		return b._db.Get(v, xquery, xargs...)
	}
	return hdsdk.Mysql.My().Get(v, xquery, xargs...)
}

func (b *SqlBuilder) Select(v any, args ...*protobuf.ListParam) error {
	var p *pagination.Pagination
	bd := b.selectBuilder
	if len(args) > 0 {
		p = pagination.NewWithParam(args[0])
		bd = bd.Offset(p.Offset).Limit(p.PageSize)
	}

	xquery, xargs, err := bd.ToSql()
	if err != nil {
		return err
	}
	if b._db != nil {
		return b._db.Select(v, xquery, xargs...)
	}

	return hdsdk.Mysql.My().Select(v, xquery, xargs...)
}

func (b *SqlBuilder) Count() (int64, error) {
	xquery, xargs, err := b.selectBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var total int64
	if b._db != nil {
		err = b._db.Select(&total, xquery, xargs...)
	} else {
		err = hdsdk.Mysql.My().Get(&total, xquery, xargs...)
	}
	return total, err
}
