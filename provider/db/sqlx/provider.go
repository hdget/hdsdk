package sqlx

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
)

type mysqlProvider struct {
	logger    intf.LoggerProvider
	config    *mysqlProviderConfig
	defaultDb intf.DbClient
	masterDb  intf.DbClient
	slaveDbs  []intf.DbClient
	extraDbs  map[string]intf.DbClient
}

const (
	// 这里设置解析时间类型https://github.com/go-sql-driver/mysql#timetime-support
	// DSN (Data Type NickName): username:password@protocol(address)/dbname?param=value
	dsnTemplate = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local"
)

func New(configProvider intf.ConfigProvider, logger intf.LoggerProvider) (intf.DbProvider, error) {
	c, err := newConfig(configProvider)
	if err != nil {
		return nil, errors.Wrap(err, "new mysql config")
	}

	provider := &mysqlProvider{
		logger:   logger,
		config:   c,
		slaveDbs: make([]intf.DbClient, len(c.Slaves)),
		extraDbs: make(map[string]intf.DbClient),
	}

	err = provider.Init(logger, c)
	if err != nil {
		logger.Fatal("init mysql provider", "err", err)
	}

	return provider, nil
}

func (p *mysqlProvider) Init(args ...any) error {
	var err error
	if p.config.Default != nil {
		p.defaultDb, err = newClient(p.config.Default)
		if err != nil {
			return errors.Wrap(err, "init mysql default connection")
		}

		p.logger.Debug("init mysql default", "host", p.config.Default.Host)
	}

	if p.config.Master != nil {
		p.masterDb, err = newClient(p.config.Master)
		if err != nil {
			return errors.Wrap(err, "init mysql master connection")
		}
		p.logger.Debug("init mysql master", "host", p.config.Master.Host)
	}

	for i, slaveConf := range p.config.Slaves {
		slaveClient, err := newClient(slaveConf)
		if err != nil {
			return errors.Wrapf(err, "init mysql slave connection, index: %d", i)
		}

		p.slaveDbs[i] = slaveClient
		p.logger.Debug("init mysql slave", "index", i, "host", slaveConf.Host)
	}

	for _, itemConf := range p.config.Items {
		itemClient, err := newClient(itemConf)
		if err != nil {
			return errors.Wrapf(err, "new mysql extra connection, name: %s", itemConf.Name)
		}

		p.extraDbs[itemConf.Name] = itemClient
		p.logger.Debug("init mysql extra", "name", itemConf.Name, "host", itemConf.Host)
	}

	return nil
}

func (p *mysqlProvider) My() intf.DbClient {
	return p.defaultDb
}

func (p *mysqlProvider) Master() intf.DbClient {
	return p.masterDb
}

func (p *mysqlProvider) Slave(i int) intf.DbClient {
	return p.slaveDbs[i]
}

func (p *mysqlProvider) By(name string) intf.DbClient {
	return p.extraDbs[name]
}
