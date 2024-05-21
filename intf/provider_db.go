package intf

import (
	"database/sql"
)

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
