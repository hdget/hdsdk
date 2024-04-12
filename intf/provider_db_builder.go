package intf

type DbBuilderProvider interface {
	Provider
	ToSql() (string, []any, error)
}
