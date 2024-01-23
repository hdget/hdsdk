package intf

import (
	"github.com/Masterminds/squirrel"
	"github.com/hdget/hdsdk/protobuf"
	"github.com/jmoiron/sqlx"
)

type DbProvider interface {
	My() DbBuilder
	Master() DbBuilder
	Slave(i int) DbBuilder
	By(name string) DbBuilder
}

type DbBuilder interface {
	Squirrel(builder squirrel.Sqlizer) ApiDbClient // squirrel builder support
	Sqrl(builder squirrel.Sqlizer) ApiDbClient     // squirrel builder support
	Db() *sqlx.DB
	ApiDbClient
}

type ApiDbClient interface {
	ToSql() (string, []any, error)
	Get(v any) error
	Select(v any, args ...*protobuf.ListParam) error
	Count() (int64, error)
	Query(args ...*protobuf.ListParam) (*sqlx.Rows, error)
}

//
//type DbClient interface {
//	SqlxClient
//	BuilderClient
//	Sq(builder squirrel.Sqlizer) DbClient // squirrel builder support
//	Sqrl(builder sqrl.Sqlizer) DbClient   // sqrl builder support
//	Db() *sqlx.DB
//}
