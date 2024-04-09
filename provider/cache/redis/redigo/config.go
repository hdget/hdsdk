package redigo

import (
	"github.com/hdget/hdsdk/v1/errdef"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/pkg/errors"
)

type providerConfig struct {
	redis *redisProviderConfig `mapstructure:"redis"` // 日志配置
}

type redisProviderConfig struct {
	Default *instanceConfig   `mapstructure:"default"`
	Items   []*instanceConfig `mapstructure:"items"`
}

type instanceConfig struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
}

func NewConfig(configProvider intf.ConfigProvider) (*redisProviderConfig, error) {
	if configProvider == nil {
		return nil, errdef.ErrInvalidConfig
	}

	var c providerConfig
	err := configProvider.UnmarshalProviderConfig(&c)
	if err != nil {
		return nil, err
	}

	err = c.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate redis provider config")
	}

	return c.redis, nil
}

func (c providerConfig) validate() error {
	if c.redis == nil {
		return errdef.ErrEmptyConfig
	}

	if c.redis.Default != nil {
		err := c.validateInstanceConfig(c.redis.Default)
		if err != nil {
			return err
		}
	}

	for _, item := range c.redis.Items {
		err := c.validateExtraInstanceConfig(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c providerConfig) validateInstanceConfig(conf *instanceConfig) error {
	if conf.Host == "" {
		return errdef.ErrInvalidConfig
	}

	// setup default config value
	if conf.Port == 0 {
		conf.Port = 6379
	}

	return nil
}

func (c providerConfig) validateExtraInstanceConfig(conf *instanceConfig) error {
	if conf.Name == "" || conf.Host == "" {
		return errdef.ErrInvalidConfig
	}

	// setup default config value
	if conf.Port == 0 {
		conf.Port = 6379
	}

	return nil
}
