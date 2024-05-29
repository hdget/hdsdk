package intf

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type DbProvider interface {
	Provider
	My() DbClient
	Master() DbClient
	Slave(i int) DbClient
	By(name string) DbClient
}

type SqlxDbProvider interface {
	Provider
	My() SqlxDbClient
	Master() SqlxDbClient
	Slave(i int) SqlxDbClient
	By(name string) SqlxDbClient
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

type SqlxDbClient interface {
	DbClient
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Beginx() (*sqlx.Tx, error)
}
