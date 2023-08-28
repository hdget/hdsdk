package types

import (
	"github.com/Masterminds/squirrel"
	"github.com/hdget/hdsdk/protobuf"
	"github.com/jmoiron/sqlx"
)

// DbProvider database能力提供
type DbProvider interface {
	Provider
	My() DbClient       // 默认我们的数据库连接
	Master() DbClient   // 主库
	Slave(int) DbClient // 指定的从库
	By(string) DbClient // 获取某个名字的数据库连接
}

type DbClient interface {
	Sq(builder squirrel.SelectBuilder) DbClient

	Select(dest any, query string, args ...any) error
	Get(dest any, query string, args ...any) error
	Queryx(query string, args ...any) (*sqlx.Rows, error)
	Rebind(query string) string

	// squirrel builder related methods
	ToSql() (string, []any, error)
	HGet(v any) error
	HSelect(v any, args ...*protobuf.ListParam) error
	HCount() (int64, error)
	HQuery(args ...*protobuf.ListParam) (*sqlx.Rows, error)
}

// database capability
const (
	_ = SdkCategoryDb + iota
	SdkTypeDbMysql
)
