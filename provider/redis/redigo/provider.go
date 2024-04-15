package redigo

import (
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
)

type redigoProvider struct {
	logger        intf.LoggerProvider
	config        *redisProviderConfig
	defaultClient intf.RedisClient            // 缺省redis
	extraClients  map[string]intf.RedisClient // 额外的redis
}

func New(configProvider intf.ConfigProvider, logger intf.LoggerProvider) (intf.RedisProvider, error) {
	c, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	provider := &redigoProvider{
		logger: logger,
		config: c,
	}

	if len(c.Items) > 0 {
		provider.extraClients = make(map[string]intf.RedisClient)
	}

	err = provider.Init()
	if err != nil {
		logger.Fatal("init redis provider", "err", err)
	}

	return provider, nil
}

func (r *redigoProvider) Init(args ...any) error {
	var err error
	if r.config.Default != nil {
		r.defaultClient, err = newRedisClient(r.config.Default)
		if err != nil {
			return errors.Wrap(err, "init redis default client")
		}
		r.logger.Debug("init redis default client", "host", r.config.Default.Host)
	}

	for _, itemConf := range r.config.Items {
		itemClient, err := newRedisClient(itemConf)
		if err != nil {
			return errors.Wrapf(err, "new redis extra client, name: %s", itemConf.Name)
		}

		r.extraClients[itemConf.Name] = itemClient
		r.logger.Debug("init redis extra client", "name", itemConf.Name, "host", itemConf.Host)
	}

	return nil
}

func (r *redigoProvider) My() intf.RedisClient {
	return r.defaultClient
}

func (r *redigoProvider) By(name string) intf.RedisClient {
	return r.extraClients[name]
}
