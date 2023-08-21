package redis

import (
	"github.com/hdget/hdsdk/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type ConfigRedis struct {
	Default *RedisConf   `mapstructure:"default"`
	Items   []*RedisConf `mapstructure:"items"`
}

type RedisConf struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
}

func validateConf(providerType string, conf *RedisConf) error {
	if conf == nil {
		return types.ErrEmptyConfig
	}

	if conf.Host == "" || conf.Port == 0 {
		return types.ErrInvalidConfig
	}

	// extra provider需要提供name
	if providerType == types.ProviderTypeOther && conf.Name == "" {
		return types.ErrInvalidConfig
	}

	return nil
}

func parseConfig(rootConfiger types.Configer) (*ConfigRedis, error) {
	if rootConfiger == nil {
		return nil, types.ErrEmptyConfig
	}

	data := rootConfiger.GetRedisConfig()
	if data == nil {
		return nil, types.ErrEmptyConfig
	}

	values, ok := data.(map[string]interface{})
	if !ok {
		return nil, types.ErrInvalidConfig
	}

	var conf ConfigRedis
	err := mapstructure.Decode(values, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "decode redis config")
	}

	return &conf, nil
}
