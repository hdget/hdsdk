package intf

import (
	"database/sql"
	"github.com/hdget/hdsdk/v2/protobuf"
	"github.com/jmoiron/sqlx"
)

type Sqlizer interface {
	ToSql() (string, []interface{}, error)
}

type DbProvider interface {
	Provider
	My() DbClient
	Master() DbClient
	Slave(i int) DbClient
	By(name string) DbClient
}

type DbClient interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Rebind(query string) string
	Close() error
}

type SqlxDbProvider interface {
	Provider
	My() SqlxDbClient
	Master() SqlxDbClient
	Slave(i int) SqlxDbClient
	By(name string) SqlxDbClient
}

type SqlxDbClient interface {
	DbClient
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Beginx() (*sqlx.Tx, error)
	Db() *sqlx.DB
}

type DbBuilderProvider interface {
	Provider
	My() DbBuilderClient
	Master() DbBuilderClient
	Slave(i int) DbBuilderClient
	By(name string) DbBuilderClient
	Set(sqlizer Sqlizer)
}

type DbBuilderClient interface {
	ToSql() (string, []any, error)
	XGet(v any) error
	XSelect(v any, args ...*protobuf.ListParam) error
	XCount() (int64, error)
	XQuery(args ...*protobuf.ListParam) (*sqlx.Rows, error)
}
