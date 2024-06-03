package sqlx_mysql

import (
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type mysqlProvider struct {
	logger    intf.LoggerProvider
	config    *mysqlProviderConfig
	defaultDb *sqlx.DB
	masterDb  *sqlx.DB
	slaveDbs  []*sqlx.DB
	extraDbs  map[string]*sqlx.DB
	_builder  intf.Sqlizer
}

func New(configProvider intf.ConfigProvider, logger intf.LoggerProvider) (intf.DbBuilderProvider, error) {
	c, err := newConfig(configProvider)
	if err != nil {
		return nil, errors.Wrap(err, "new mysql config")
	}

	provider := &mysqlProvider{
		logger:   logger,
		config:   c,
		slaveDbs: make([]*sqlx.DB, len(c.Slaves)),
		extraDbs: make(map[string]*sqlx.DB),
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
		p.defaultDb, err = newDb(p.config.Default)
		if err != nil {
			return errors.Wrap(err, "init mysql default connection")
		}

		p.logger.Debug("init mysql default", "host", p.config.Default.Host)
	}

	if p.config.Master != nil {
		p.masterDb, err = newDb(p.config.Master)
		if err != nil {
			return errors.Wrap(err, "init mysql master connection")
		}
		p.logger.Debug("init mysql master", "host", p.config.Master.Host)
	}

	for i, slaveConf := range p.config.Slaves {
		slaveClient, err := newDb(slaveConf)
		if err != nil {
			return errors.Wrapf(err, "init mysql slave connection, index: %d", i)
		}

		p.slaveDbs[i] = slaveClient
		p.logger.Debug("init mysql slave", "index", i, "host", slaveConf.Host)
	}

	for _, itemConf := range p.config.Items {
		itemClient, err := newDb(itemConf)
		if err != nil {
			return errors.Wrapf(err, "new mysql extra connection, name: %s", itemConf.Name)
		}

		p.extraDbs[itemConf.Name] = itemClient
		p.logger.Debug("init mysql extra", "name", itemConf.Name, "host", itemConf.Host)
	}

	return nil
}

func (p *mysqlProvider) My() intf.DbBuilderClient {
	return &mysqlClient{
		DB:       p.defaultDb,
		_builder: p._builder,
	}
}

func (p *mysqlProvider) Master() intf.DbBuilderClient {
	return &mysqlClient{
		DB:       p.masterDb,
		_builder: p._builder,
	}
}

func (p *mysqlProvider) Slave(i int) intf.DbBuilderClient {
	return &mysqlClient{
		DB:       p.slaveDbs[i],
		_builder: p._builder,
	}
}

func (p *mysqlProvider) By(name string) intf.DbBuilderClient {
	return &mysqlClient{
		DB:       p.extraDbs[name],
		_builder: p._builder,
	}
}

func (p *mysqlProvider) Set(builder intf.Sqlizer) {
	p._builder = builder
}
