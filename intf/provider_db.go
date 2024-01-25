package intf

import (
	"database/sql"
)

type DbProvider interface {
	My() DbClient
	Master() DbClient
	Slave(i int) DbClient
	By(name string) DbClient
}

type DbClient interface {
	UseBuilder(builder DbBuilder)
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

type DbApiExtension interface {
	NamedExec(query string, arg any) (sql.Result, error) // 命名
	One(v any) error
	All(v any) error
	Count() (int64, error)
}
