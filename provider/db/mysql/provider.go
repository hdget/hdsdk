package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

	return provider, nil
}

func (p *mysqlProvider) Init(args ...any) error {
	var err error
	if p.config.Default != nil {
		p.defaultDb, err = newDB(p.config.Default)
		if err != nil {
			return errors.Wrap(err, "init mysql default connection")
		}

		// 设置boil的缺省db
		boil.SetDB(p.defaultDb)
		p.logger.Debug("init mysql default", "host", p.config.Default.Host)
	}

	if p.config.Master != nil {
		p.masterDb, err = newDB(p.config.Master)
		if err != nil {
			return errors.Wrap(err, "init mysql master connection")
		}
		p.logger.Debug("init mysql master", "host", p.config.Master.Host)
	}

	for i, slaveConf := range p.config.Slaves {
		slaveClient, err := newDB(slaveConf)
		if err != nil {
			return errors.Wrapf(err, "init mysql slave connection, index: %d", i)
		}

		p.slaveDbs[i] = slaveClient
		p.logger.Debug("init mysql slave", "index", i, "host", slaveConf.Host)
	}

	for _, itemConf := range p.config.Items {
		itemClient, err := newDB(itemConf)
		if err != nil {
			return errors.Wrapf(err, "new mysql extra connection, name: %s", itemConf.Name)
		}

		p.extraDbs[itemConf.Name] = itemClient
		p.logger.Debug("init mysql extra", "name", itemConf.Name, "host", itemConf.Host)
	}

	return nil
}

func (p *mysqlProvider) My() intf.DbClient {
	return newClient(p.defaultDb)
}

func (p *mysqlProvider) Master() intf.DbClient {
	return newClient(p.masterDb)
}

func (p *mysqlProvider) Slave(i int) intf.DbClient {
	var db *sqlx.DB
	if i <= len(p.slaveDbs) {
		db = p.slaveDbs[i]
	}
	return newClient(db)
}

func (p *mysqlProvider) By(name string) intf.DbClient {
	return newClient(p.extraDbs[name])
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
