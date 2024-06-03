package sqlx_mysql

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdsdk/v2/lib/pagination"
	"github.com/hdget/hdsdk/v2/protobuf"
	"github.com/hdget/hdutils/page"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type mysqlClient struct {
	*sqlx.DB
	_builder squirrel.Sqlizer
}

const (
	// 这里设置解析时间类型https://github.com/go-sql-driver/mysql#timetime-support
	// DSN (Data Type NickName): username:password@protocol(address)/dbname?param=value
	dsnTemplate = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local"
)

func newClient(c *mysqlConfig) (intf.DbBuilderClient, error) {
	// 构造连接参数
	dsn := fmt.Sprintf(dsnTemplate, c.User, c.Password, c.Host, c.Port, c.Database)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// https://www.alexedwards.net/blog/configuring-sqldb
	// https://making.pusher.com/production-ready-connection-pooling-in-go
	// Avoid issue:
	// packets.go:123: closing bad idle connection: EOF
	// connection.go:173: driver: bad connection
	db.SetConnMaxLifetime(3 * time.Minute)

	return &mysqlClient{DB: db}, nil
}

func (m mysqlClient) ToSql() (string, []any, error) {
	return m._builder.ToSql()
}

func (m mysqlClient) XGet(v any) error {
	xquery, xargs, err := m._builder.ToSql()
	if err != nil {
		return err
	}

	return m.DB.Get(v, xquery, xargs...)
}

func (m mysqlClient) XSelect(v any, args ...*protobuf.ListParam) error {
	bd := m._builder
	if len(args) > 0 {
		selBd, ok := bd.(squirrel.SelectBuilder)
		if !ok {
			return errors.New("invalid select builder")
		}

		p := pagination.New(args[0])
		bd = selBd.Offset(p.Offset).Limit(p.PageSize)
	}

	xquery, xargs, err := bd.ToSql()
	if err != nil {
		return err
	}
	return m.DB.Select(v, xquery, xargs...)
}

func (m mysqlClient) XCount() (int64, error) {
	xquery, xargs, err := m._builder.ToSql()
	if err != nil {
		return 0, err
	}

	var total int64
	err = m.DB.Get(&total, xquery, xargs...)
	return total, err
}

func (m mysqlClient) XQuery(args ...*protobuf.ListParam) (*sqlx.Rows, error) {
	var p page.Pagination
	bd := m._builder
	if len(args) > 0 {
		selBd, ok := bd.(squirrel.SelectBuilder)
		if !ok {
			return nil, errors.New("invalid select builder")
		}

		p = pagination.New(args[0])
		bd = selBd.Offset(p.Offset).Limit(p.PageSize)
	}

	xquery, xargs, err := bd.ToSql()
	if err != nil {
		return nil, err
	}
	return m.DB.Queryx(xquery, xargs...)
}
