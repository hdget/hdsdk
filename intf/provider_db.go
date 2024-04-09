package intf

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type DbProvider interface {
	My() DbClient
	Master() DbClient
	Slave(i int) DbClient
	By(name string) DbClient
}

type DbClient interface {
	UseBuilder(builder DbBuilder) DbClient
	DbApiBasic
	DbApiExtension
}

type DbBuilder interface {
	ToSql() (string, []any, error)
}

type DbApiBasic interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type DbApiSqlx interface {
	sqlx.Ext
	sqlx.ExtContext
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

type DbApiExtension interface {
	DbApiSqlx
	One(v any) error
	All(v any) error
	Count() (int64, error)
}
