package types

import "github.com/jmoiron/sqlx"

// DbProvider database能力提供
type DbProvider interface {
	Provider

	My() *sqlx.DB // 默认我们的数据库连接
	Master() *sqlx.DB
	Slave(int) *sqlx.DB
	By(string) *sqlx.DB // 获取某个名字的数据库连接
}

type DbClient interface {
	// Rebind transforms a query from QUESTION to the DB driver's bindvar type.
	Rebind(query string) string

	Get(dest interface{}, query string, args ...interface{}) error
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
}

// database capability
const (
	_ = SdkCategoryDb + iota
	SdkTypeDbMysql
)
