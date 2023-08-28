package db

import (
	"github.com/Masterminds/squirrel"
	"github.com/hdget/hdsdk/lib/pagination"
	"github.com/hdget/hdsdk/protobuf"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type SqurrelClient struct {
	*BaseDbClient
	_builder squirrel.Sqlizer
}

func (b *SqurrelClient) XGet(v any) error {
	xquery, xargs, err := b._builder.ToSql()
	if err != nil {
		return err
	}

	return b.DB.Get(v, xquery, xargs...)
}

func (b *SqurrelClient) XSelect(v any, args ...*protobuf.ListParam) error {
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
	return b.DB.Select(v, xquery, xargs...)
}

func (b *SqurrelClient) XCount() (int64, error) {
	xquery, xargs, err := b._builder.ToSql()
	if err != nil {
		return 0, err
	}

	var total int64
	err = b.DB.Get(&total, xquery, xargs...)
	return total, err
}

func (b *SqurrelClient) XQuery(args ...*protobuf.ListParam) (*sqlx.Rows, error) {
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
	return b.DB.Queryx(xquery, xargs...)
}

func (b *SqurrelClient) ToSql() (string, []any, error) {
	return b._builder.ToSql()
}
