package db

import (
	"github.com/elgris/sqrl"
	"github.com/hdget/hdsdk/lib/pagination"
	"github.com/hdget/hdsdk/protobuf"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type SqrlDbClient struct {
	*sqlx.DB
	_builder sqrl.Sqlizer
}

func (b *SqrlDbClient) HGet(v any) error {
	xquery, xargs, err := b._builder.ToSql()
	if err != nil {
		return err
	}

	return b.DB.Get(v, xquery, xargs...)
}

func (b *SqrlDbClient) HSelect(v any, args ...*protobuf.ListParam) error {
	var p *pagination.Pagination
	bd := b._builder
	if len(args) > 0 {
		selBd, ok := bd.(*sqrl.SelectBuilder)
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

func (b *SqrlDbClient) HCount() (int64, error) {
	xquery, xargs, err := b._builder.ToSql()
	if err != nil {
		return 0, err
	}

	var total int64
	err = b.DB.Get(&total, xquery, xargs...)
	return total, err
}

func (b *SqrlDbClient) HQuery(args ...*protobuf.ListParam) (*sqlx.Rows, error) {
	var p *pagination.Pagination
	bd := b._builder
	if len(args) > 0 {
		selBd, ok := bd.(*sqrl.SelectBuilder)
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

func (b *SqrlDbClient) ToSql() (string, []any, error) {
	return b._builder.ToSql()
}
