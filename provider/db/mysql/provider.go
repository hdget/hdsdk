package mysql

import (
	"github.com/hdget/hdsdk/intf"
	"github.com/pkg/errors"
)

type mysqlProvider struct {
	defaultClient intf.DbClient
	masterClient  intf.DbClient
	slaveClients  []intf.DbClient
	extraClients  map[string]intf.DbClient
}

func New(mysqlProviderConfig *mysqlProviderConfig, logger intf.Logger) (intf.DbProvider, error) {
	provider := &mysqlProvider{
		slaveClients: make([]intf.DbClient, len(mysqlProviderConfig.Slaves)),
		extraClients: make(map[string]intf.DbClient),
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
		m.defaultClient, err = newMysqlClient(providerConfig.Default)
		if err != nil {
			return errors.Wrap(err, "new mysql default client")
		}
	}

	if providerConfig.Master != nil {
		m.masterClient, err = newMysqlClient(providerConfig.Master)
		if err != nil {
			return errors.Wrap(err, "new mysql master client")
		}
	}

	for i, slaveConf := range providerConfig.Slaves {
		slaveClient, err := newMysqlClient(slaveConf)
		if err != nil {
			return errors.Wrapf(err, "new mysql slave client, index: %d", i)
		}

		m.slaveClients[i] = slaveClient
	}

	for _, itemConf := range providerConfig.Items {
		itemClient, err := newMysqlClient(itemConf)
		if err != nil {
			return errors.Wrapf(err, "new mysql extra client, name: %d", itemConf.Name)
		}

		m.extraClients[itemConf.Name] = itemClient
	}

	return nil
}

func (m mysqlProvider) My() intf.DbClient {
	return m.defaultClient
}

func (m mysqlProvider) Master() intf.DbClient {
	return m.masterClient
}

func (m mysqlProvider) Slave(i int) intf.DbClient {
	if i >= len(m.slaveClients) {
		return nil
	}
	return m.slaveClients[i]
}

func (m mysqlProvider) By(name string) intf.DbClient {
	return m.extraClients[name]
}
