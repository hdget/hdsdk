package mysql

import (
	"fmt"
	"github.com/hdget/hdsdk/intf"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type mysqlClient struct {
	DB *sqlx.DB
}

func (c *mysqlClient) Get(v any) error {
	//TODO implement me
	panic("implement me")
}

func (c *mysqlClient) Select(v any, args ...*interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *mysqlClient) Count() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (c *mysqlClient) Query(args ...*interface{}) (*sqlx.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func newMysqlClient(mysqlConfig *mysqlConfig) (intf.DbClient, error) {
	client := &mysqlClient{}
	err := client.connect(mysqlConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *mysqlClient) connect(mysqlConfig *mysqlConfig) error {
	// 这里设置解析时间类型https://github.com/go-sql-driver/mysql#timetime-support
	// DSN (Data Type NickName): username:password@protocol(address)/dbname?param=value
	t := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local"
	// 构造连接参数
	connStr := fmt.Sprintf(t, mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Database)
	instance, err := sqlx.Connect("mysql", connStr)
	if err != nil {
		return errors.Wrap(err, "mysql connect")
	}

	// https://www.alexedwards.net/blog/configuring-sqldb
	// https://making.pusher.com/production-ready-connection-pooling-in-go
	// Avoid issue:
	// packets.go:123: closing bad idle connection: EOF
	// connection.go:173: driver: bad connection
	instance.SetConnMaxLifetime(3 * time.Minute)

	c.DB = instance
	return nil
}
