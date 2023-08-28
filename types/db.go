package types

import (
	"database/sql"
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
	DriverName() string
	MapperFunc(mf func(string) string)
	Rebind(query string) string
	Unsafe() *sqlx.DB
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	MustBegin() *sqlx.Tx
	Beginx() (*sqlx.Tx, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	MustExec(query string, args ...interface{}) sql.Result
	Preparex(query string) (*sqlx.Stmt, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)

	// squirrel builder related methods
	Sq(builder squirrel.SelectBuilder) DbClient
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
