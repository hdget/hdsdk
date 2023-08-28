package db

import (
	"github.com/Masterminds/squirrel"
	"github.com/elgris/sqrl"
	"github.com/hdget/hdsdk/protobuf"
	"github.com/hdget/hdsdk/types"
	"github.com/jmoiron/sqlx"
)

type BaseDbClient struct {
	*sqlx.DB
}

var (
	_ types.DbClient = (*BaseDbClient)(nil)
)

func (b *BaseDbClient) Sq(builder squirrel.Sqlizer) types.DbClient {
	return &SqurrelClient{
		BaseDbClient: b,
		_builder:     builder,
	}
}

func (b *BaseDbClient) Sqrl(builder sqrl.Sqlizer) types.DbClient {
	return &SqrlClient{
		BaseDbClient: b,
		_builder:     builder,
	}
}

func (b *BaseDbClient) ToSql() (string, []any, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BaseDbClient) XGet(v any) error {
	//TODO implement me
	panic("implement me")
}

func (b *BaseDbClient) XSelect(v any, args ...*protobuf.ListParam) error {
	//TODO implement me
	panic("implement me")
}

func (b *BaseDbClient) XCount() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BaseDbClient) XQuery(args ...*protobuf.ListParam) (*sqlx.Rows, error) {
	//TODO implement me
	panic("implement me")
}
