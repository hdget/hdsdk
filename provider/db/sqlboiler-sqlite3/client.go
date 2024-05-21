package sqlboiler_sqlite3

import (
	"database/sql"
	"fmt"
	"github.com/hdget/hdsdk/v2/intf"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

type sqliteClient struct {
	*sql.DB
}

const (
	// 这里设置解析时间类型https://github.com/go-sql-driver/mysql#timetime-support
	// DSN (Data Type NickName): username:password@protocol(address)/dbname?param=value
	dsnTemplate = "file:%s?_loc=Local"
)

func newClient(c *sqliteProviderConfig) (intf.DbClient, error) {
	instance, err := newInstance(c)
	if err != nil {
		return nil, err
	}
	return &sqliteClient{DB: instance}, nil
}

func newInstance(c *sqliteProviderConfig) (*sql.DB, error) {
	// 构造连接参数
	dsn := fmt.Sprintf(dsnTemplate, c.DbName)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	var userVersion int
	err = db.QueryRow("PRAGMA user_version").Scan(&userVersion)
	if err != nil {
		_ = db.Close()

		dir, _ := os.Getwd()
		return nil, fmt.Errorf("fail connect db: %s", filepath.Join(dir, c.DbName))
	}

	// https://www.alexedwards.net/blog/configuring-sqldb
	// https://making.pusher.com/production-ready-connection-pooling-in-go
	// Avoid issue:
	// packets.go:123: closing bad idle connection: EOF
	// connection.go:173: driver: bad connection
	db.SetConnMaxLifetime(3 * time.Minute)
	return db, nil
}

func (m sqliteClient) Get(dest interface{}, query string, args ...interface{}) error {
	//TODO implement me
	return nil
}

func (m sqliteClient) Select(dest interface{}, query string, args ...interface{}) error {
	//TODO implement me
	return nil
}

func (m sqliteClient) Rebind(query string) string {
	return ""
}
