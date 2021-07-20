// @Title  log capability of zerolog
// @Description  zerolog implementation of log capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package redis

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/hdget/sdk/lib/cache/provider"
	"github.com/hdget/sdk/types"
)

type RedisProvider struct {
	provider.BaseCacheProvider
	Log types.LogProvider
}

var (
	_ types.Provider      = (*RedisProvider)(nil)
	_ types.CacheProvider = (*RedisProvider)(nil)
)

// @desc	implements types.Provider interface, used to initialize the capability
// @author	Ryan Fan	(2021-06-09)
// @param	baseconf.Configer	root config interface to extract config info
// @return	error
func (rp *RedisProvider) Init(rootConfiger types.Configer, logger types.LogProvider, args ...interface{}) error {
	// 获取日志配置信息
	config, err := parseConfig(rootConfiger)
	if err != nil {
		return err
	}

	// 缺省redis必须要配置合法
	err = validateConf(types.PROVIDER_TYPE_DEFAULT, config.Default)
	if err != nil {
		logger.Fatal("validate redis config", "type", types.PROVIDER_TYPE_DEFAULT, "err", err)
	}

	rp.Default, err = rp.connect(config.Default)
	if err != nil {
		logger.Fatal("connect redis", "type", types.PROVIDER_TYPE_DEFAULT, "host", config.Default.Host, "err", err)
	}
	logger.Debug("connect redis", "type", types.PROVIDER_TYPE_DEFAULT, "host", config.Default.Host)

	// 额外的redis
	rp.Items = make(map[string]types.CacheClient)
	for _, otherConf := range config.Items {
		if err := validateConf(types.PROVIDER_TYPE_OTHER, otherConf); err == nil {
			instance, err := rp.connect(otherConf)
			if instance != nil {
				rp.Items[otherConf.Name] = instance
			}
			logger.Debug("connect redis", "type", otherConf.Name, "host", otherConf.Host, "err", err)
		}
	}

	return nil
}

func (rp *RedisProvider) connect(conf *RedisConf) (types.CacheClient, error) {
	client := NewRedisClient(conf)
	err := client.Ping()
	if err != nil {
		return nil, err
	}
	return client, nil
}
