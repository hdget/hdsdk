package redis

import (
	"github.com/hdget/hdsdk/intf"
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
		return intf.ErrEmptyConfig
	}

	if conf.Host == "" || conf.Port == 0 {
		return intf.ErrInvalidConfig
	}

	// extra provider需要提供name
	if providerType == intf.ProviderTypeOther && conf.Name == "" {
		return intf.ErrInvalidConfig
	}

	return nil
}

func parseConfig(rootConfiger intf.SdkConfig) (*ConfigRedis, error) {
	if rootConfiger == nil {
		return nil, intf.ErrEmptyConfig
	}

	data := rootConfiger.GetRedisConfig()
	if data == nil {
		return nil, intf.ErrEmptyConfig
	}

	values, ok := data.(map[string]interface{})
	if !ok {
		return nil, intf.ErrInvalidConfig
	}

	var conf ConfigRedis
	err := mapstructure.Decode(values, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "decode redis configer")
	}

	return &conf, nil
}
