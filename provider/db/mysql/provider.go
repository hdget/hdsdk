package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hdget/hdsdk/v1/instance"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type mysqlProvider struct {
	logger    intf.LoggerProvider
	config    *mysqlProviderConfig
	defaultDb *sqlx.DB
	masterDb  *sqlx.DB
	slaveDbs  []*sqlx.DB
	extraDbs  map[string]*sqlx.DB
}

func New(mysqlConfig *mysqlProviderConfig, logger intf.LoggerProvider) (intf.DbProvider, error) {
	provider := &mysqlProvider{
		logger:   logger,
		config:   mysqlConfig,
		slaveDbs: make([]*sqlx.DB, len(mysqlConfig.Slaves)),
		extraDbs: make(map[string]*sqlx.DB),
	}

	err := provider.Init(logger, mysqlConfig)
	if err != nil {
		logger.Fatal("init mysql provider", "err", err)
	}

	instance.Register(provider)

	return provider, nil
}

func (m *mysqlProvider) Init(args ...any) error {
	if len(args) == 0 {
		return errors.New("need mysql provider config")
	}

	c, ok := args[0].(*mysqlProviderConfig)
	if !ok {
		return errors.New("invalid mysql provider config")
	}

	var err error
	if c.Default != nil {
		m.defaultDb, err = newDB(c.Default)
		if err != nil {
			return errors.Wrap(err, "init mysql default connection")
		}
		m.logger.Debug("init mysql default", "host", c.Default.Host)
	}

	if c.Master != nil {
		m.masterDb, err = newDB(c.Master)
		if err != nil {
			return errors.Wrap(err, "init mysql master connection")
		}
		m.logger.Debug("init mysql master", "host", c.Master.Host)
	}

	for i, slaveConf := range c.Slaves {
		slaveClient, err := newDB(slaveConf)
		if err != nil {
			return errors.Wrapf(err, "init mysql slave connection, index: %d", i)
		}

		m.slaveDbs[i] = slaveClient
		m.logger.Debug("init mysql slave", "index", i, "host", slaveConf.Host)
	}

	for _, itemConf := range c.Items {
		itemClient, err := newDB(itemConf)
		if err != nil {
			return errors.Wrapf(err, "new mysql extra connection, name: %s", itemConf.Name)
		}

		m.extraDbs[itemConf.Name] = itemClient
		m.logger.Debug("init mysql extra", "name", itemConf.Name, "host", itemConf.Host)
	}

	return nil
}

func (m *mysqlProvider) My() intf.DbClient {
	return newClient(m.defaultDb)
}

func (m *mysqlProvider) Master() intf.DbClient {
	return newClient(m.masterDb)
}

func (m *mysqlProvider) Slave(i int) intf.DbClient {
	var db *sqlx.DB
	if i <= len(m.slaveDbs) {
		db = m.slaveDbs[i]
	}
	return newClient(db)
}

func (m *mysqlProvider) By(name string) intf.DbClient {
	return newClient(m.extraDbs[name])
}

func newDB(mysqlConfig *instanceConfig) (*sqlx.DB, error) {
	// 这里设置解析时间类型https://github.com/go-hdsql-driver/mysql#timetime-support
	// DSN (Data Type NickName): username:password@protocol(address)/dbname?param=value
	t := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local"
	// 构造连接参数
	connStr := fmt.Sprintf(t, mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Database)
	instance, err := sqlx.Connect("mysql", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "mysql connect")
	}

	// https://www.alexedwards.net/blog/configuring-sqldb
	// https://making.pusher.com/production-ready-connection-pooling-in-go
	// Avoid issue:
	// packets.go:123: closing bad idle connection: EOF
	// connection.go:173: driver: bad connection
	instance.SetConnMaxLifetime(3 * time.Minute)
	return instance, nil
}
