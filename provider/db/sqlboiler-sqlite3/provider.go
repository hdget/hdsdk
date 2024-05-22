package sqlboiler_sqlite3

import (
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	_ "modernc.org/sqlite"
)

type sqliteProvider struct {
	logger    intf.LoggerProvider
	config    *sqliteProviderConfig
	defaultDb intf.DbClient
}

func New(configProvider intf.ConfigProvider, logger intf.LoggerProvider) (intf.DbProvider, error) {
	c, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	provider := &sqliteProvider{
		logger: logger,
		config: c,
	}

	err = provider.Init(logger, c)
	if err != nil {
		logger.Fatal("init sqlite3 provider", "err", err)
	}

	return provider, nil
}

func (p *sqliteProvider) Init(args ...any) error {
	var err error

	p.defaultDb, err = newClient(p.config)
	if err != nil {
		return errors.Wrap(err, "new sqlite3 client")
	}

	// 设置boil的缺省db
	boil.SetDB(p.defaultDb)
	p.logger.Debug("init sqlite3 provider", "db", p.config.DbName)

	return nil
}

// Connect 从指定的文件创建创建数据库连接
func Connect(dbFile string) (intf.DbClient, error) {
	client, err := newClient(nil, dbFile)
	if err != nil {
		return nil, errors.Wrapf(err, "connect sqlite3: %s", dbFile)
	}

	// 设置boil的缺省db
	boil.SetDB(client)
	return client, nil
}

func (p *sqliteProvider) My() intf.DbClient {
	return p.defaultDb
}

func (p *sqliteProvider) Master() intf.DbClient {
	return nil
}

func (p *sqliteProvider) Slave(i int) intf.DbClient {
	return nil
}

func (p *sqliteProvider) By(name string) intf.DbClient {
	return nil
}
