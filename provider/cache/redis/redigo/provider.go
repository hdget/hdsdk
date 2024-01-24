package redigo

import (
	"github.com/hdget/hdsdk/v1/errdef"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/pkg/errors"
)

type redigoProvider struct {
	defaultClient intf.RedisClient            // 缺省redis
	extraClients  map[string]intf.RedisClient // 额外的redis
}

func New(providerConfig *redisProviderConfig, logger intf.Logger) (intf.RedisProvider, error) {
	if providerConfig == nil {
		return nil, errdef.ErrEmptyConfig
	}

	provider := &redigoProvider{}
	if len(providerConfig.Items) > 0 {
		provider.extraClients = make(map[string]intf.RedisClient)
	}

	err := provider.Init(logger, providerConfig)
	if err != nil {
		logger.Fatal("init mysql provider", "err", err)
	}

	return provider, nil
}

func (r *redigoProvider) Init(logger intf.Logger, args ...any) error {
	if len(args) == 0 {
		return errors.New("need redis provider config")
	}

	providerConfig, ok := args[0].(*redisProviderConfig)
	if !ok {
		return errors.New("invalid redis provider config")
	}

	var err error
	if providerConfig.Default != nil {
		r.defaultClient, err = newRedisClient(providerConfig.Default)
		if err != nil {
			return errors.Wrap(err, "init redis default client")
		}
		logger.Debug("init redis default", "host", providerConfig.Default.Host)
	}

	for _, itemConf := range providerConfig.Items {
		itemClient, err := newRedisClient(itemConf)
		if err != nil {
			return errors.Wrapf(err, "new redis extra client, name: %s", itemConf.Name)
		}

		r.extraClients[itemConf.Name] = itemClient
		logger.Debug("init redis extra client", "name", itemConf.Name, "host", itemConf.Host)
	}

	return nil
}

func (r *redigoProvider) My() intf.RedisClient {
	return r.defaultClient
}

func (r *redigoProvider) By(name string) intf.RedisClient {
	return r.extraClients[name]
}
