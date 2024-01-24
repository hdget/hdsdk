package redigo

import (
	"github.com/hdget/hdsdk/v1/errdef"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type redisProviderConfig struct {
	Default *redisConfig   `mapstructure:"default"`
	Items   []*redisConfig `mapstructure:"items"`
}

type redisConfig struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
}

func NewConfig(sdkConfiger intf.SdkConfiger) (*redisProviderConfig, error) {
	if sdkConfiger == nil {
		return nil, errdef.ErrInvalidConfig
	}

	// if logger config not found, do nothing
	values := sdkConfiger.GetRedisConfig()
	if len(values) == 0 {
		return nil, errdef.ErrInvalidConfig
	}

	var providerConfig redisProviderConfig
	err := mapstructure.Decode(values, &providerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "decode redis provider config")
	}

	err = providerConfig.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate redis provider config")
	}

	return &providerConfig, nil
}

func (rp *redisProviderConfig) validate() error {
	if rp.Default != nil {
		err := rp.validateRedisConfig(rp.Default)
		if err != nil {
			return err
		}
	}

	for _, item := range rp.Items {
		err := rp.validateExtraRedisConfig(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rp *redisProviderConfig) validateRedisConfig(conf *redisConfig) error {
	if conf.Host == "" {
		return intf.ErrInvalidConfig
	}

	// setup default config value
	if conf.Port == 0 {
		conf.Port = 6379
	}

	return nil
}

func (rp *redisProviderConfig) validateExtraRedisConfig(conf *redisConfig) error {
	if conf.Name == "" || conf.Host == "" {
		return intf.ErrInvalidConfig
	}

	// setup default config value
	if conf.Port == 0 {
		conf.Port = 6379
	}

	return nil
}
