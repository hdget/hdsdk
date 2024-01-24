package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hdget/hdsdk/intf"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type mysqlProvider struct {
	builder   intf.DbBuilder
	defaultDb *sqlx.DB
	masterDb  *sqlx.DB
	slaveDbs  []*sqlx.DB
	extraDbs  map[string]*sqlx.DB
}

func New(mysqlProviderConfig *mysqlProviderConfig, logger intf.Logger) (intf.DbProvider, error) {
	provider := &mysqlProvider{
		slaveDbs: make([]*sqlx.DB, len(mysqlProviderConfig.Slaves)),
		extraDbs: make(map[string]*sqlx.DB),
	}

	err := provider.Init(logger, mysqlProviderConfig)
	if err != nil {
		logger.Fatal("init mysql provider", "err", err)
	}

	return provider, nil
}

func (m *mysqlProvider) Init(logger intf.Logger, args ...any) error {
	if len(args) == 0 {
		return errors.New("need mysql provider config")
	}

	providerConfig, ok := args[0].(*mysqlProviderConfig)
	if !ok {
		return errors.New("invalid mysql provider config")
	}

	var err error
	if providerConfig.Default != nil {
		m.defaultDb, err = newDB(providerConfig.Default)
		if err != nil {
			return errors.Wrap(err, "new mysql default connection")
		}
		logger.Debug("initialize mysql default", "host", providerConfig.Default.Host)
	}

	if providerConfig.Master != nil {
		m.masterDb, err = newDB(providerConfig.Master)
		if err != nil {
			return errors.Wrap(err, "new mysql master connection")
		}
		logger.Debug("initialize mysql master", "host", providerConfig.Master.Host)
	}

	for i, slaveConf := range providerConfig.Slaves {
		slaveClient, err := newDB(slaveConf)
		if err != nil {
			return errors.Wrapf(err, "new mysql slave connection, index: %d", i)
		}

		m.slaveDbs[i] = slaveClient
		logger.Debug("initialize mysql slave", "index", i, "host", slaveConf.Host)
	}

	for _, itemConf := range providerConfig.Items {
		itemClient, err := newDB(itemConf)
		if err != nil {
			return errors.Wrapf(err, "new mysql extra connection, name: %d", itemConf.Name)
		}

		m.extraDbs[itemConf.Name] = itemClient
		logger.Debug("initialize mysql extra", "name", itemConf.Name, "host", itemConf.Host)
	}

	return nil
}

func (m mysqlProvider) My() intf.DbClient {
	return newClient(m.defaultDb)
}

func (m mysqlProvider) Master() intf.DbClient {
	return newClient(m.masterDb)
}

func (m mysqlProvider) Slave(i int) intf.DbClient {
	var db *sqlx.DB
	if i <= len(m.slaveDbs) {
		db = m.slaveDbs[i]
	}
	return newClient(db)
}

func (m mysqlProvider) By(name string) intf.DbClient {
	return newClient(m.extraDbs[name])
}

func newDB(mysqlConfig *mysqlConfig) (*sqlx.DB, error) {
	// 这里设置解析时间类型https://github.com/go-sql-driver/mysql#timetime-support
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
