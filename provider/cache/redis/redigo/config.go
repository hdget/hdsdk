package redigo

import (
	"github.com/hdget/hdsdk/v2/errdef"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
)

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

const (
	configSection = "sdk.redis"
)

func NewConfig(configLoader intf.ConfigLoader) (*redisProviderConfig, error) {
	if configLoader == nil {
		return nil, errdef.ErrInvalidConfig
	}

	var c *redisProviderConfig
	err := configLoader.Unmarshal(&c, configSection)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, errdef.ErrEmptyConfig
	}

	err = c.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate redis provider config")
	}

	return c, nil
}

func (c *redisProviderConfig) validate() error {
	if c.Default != nil {
		err := c.validateInstanceConfig(c.Default)
		if err != nil {
			return err
		}
	}

	for _, item := range c.Items {
		err := c.validateExtraInstanceConfig(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *redisProviderConfig) validateInstanceConfig(conf *instanceConfig) error {
	if conf.Host == "" {
		return errdef.ErrInvalidConfig
	}

	// setup default config value
	if conf.Port == 0 {
		conf.Port = 6379
	}

	return nil
}

func (c *redisProviderConfig) validateExtraInstanceConfig(conf *instanceConfig) error {
	if conf.Name == "" || conf.Host == "" {
		return errdef.ErrInvalidConfig
	}

	// setup default config value
	if conf.Port == 0 {
		conf.Port = 6379
	}

	return nil
}
