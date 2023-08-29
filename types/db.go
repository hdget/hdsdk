package types

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/elgris/sqrl"
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

type SqlxClient interface {
	sqlx.Ext
	sqlx.ExtContext
	MapperFunc(mf func(string) string)
	Unsafe() *sqlx.DB
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	MustBegin() *sqlx.Tx
	Beginx() (*sqlx.Tx, error)
	MustExec(query string, args ...interface{}) sql.Result
	Preparex(query string) (*sqlx.Stmt, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}

type BuilderClient interface {
	ToSql() (string, []any, error)
	XGet(v any) error
	XSelect(v any, args ...*protobuf.ListParam) error
	XCount() (int64, error)
	XQuery(args ...*protobuf.ListParam) (*sqlx.Rows, error)
}

type DbClient interface {
	SqlxClient
	BuilderClient
	Sq(builder squirrel.Sqlizer) DbClient // squirrel builder support
	Sqrl(builder sqrl.Sqlizer) DbClient   // sqrl builder support
	Db() *sqlx.DB
}

// database capability
const (
	_ = SdkCategoryDb + iota
	SdkTypeDbMysql
)
