package mysql

import (
	"database/sql"
	"fmt"
	"github.com/hdget/hdsdk/v2/intf"
	"time"
)

type mysqlClient struct {
	*sql.DB
}

func newClient(c *mysqlConfig) (intf.DbClient, error) {
	instance, err := newInstance(c)
	if err != nil {
		return nil, err
	}
	return &mysqlClient{DB: instance}, nil
}

func newInstance(c *mysqlConfig) (*sql.DB, error) {
	// 构造连接参数
	dsn := fmt.Sprintf(dsnTemplate, c.User, c.Password, c.Host, c.Port, c.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	// https://www.alexedwards.net/blog/configuring-sqldb
	// https://making.pusher.com/production-ready-connection-pooling-in-go
	// Avoid issue:
	// packets.go:123: closing bad idle connection: EOF
	// connection.go:173: driver: bad connection
	db.SetConnMaxLifetime(3 * time.Minute)
	return db, nil
}
