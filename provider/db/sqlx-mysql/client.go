package sqlx_mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/jmoiron/sqlx"
	"time"
)

type mysqlClient struct {
	*sqlx.DB
}

const (
	// 这里设置解析时间类型https://github.com/go-sql-driver/mysql#timetime-support
	// DSN (Data Type NickName): username:password@protocol(address)/dbname?param=value
	dsnTemplate = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local"
)

func newClient(c *mysqlConfig) (intf.DbClient, error) {
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

func (m mysqlClient) Close() error {
	return m.DB.Close()
}
