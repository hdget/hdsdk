package mysql

import (
	"github.com/hdget/hdsdk/errdef"
	"github.com/hdget/hdsdk/intf"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type mysqlProviderConfig struct {
	Default *mysqlConfig   `mapstructure:"default"`
	Master  *mysqlConfig   `mapstructure:"master"`
	Slaves  []*mysqlConfig `mapstructure:"slaves"`
	Items   []*mysqlConfig `mapstructure:"items"`
}

type mysqlConfig struct {
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Timeout  int    `mapstructure:"timeout"`
}

func NewConfig(sdkConfiger intf.SdkConfiger) (*mysqlProviderConfig, error) {
	if sdkConfiger == nil {
		return nil, errdef.ErrInvalidConfig
	}

	// if logger config not found, do nothing
	values := sdkConfiger.GetMysqlConfig()
	if len(values) == 0 {
		return nil, errdef.ErrInvalidConfig
	}

	var conf mysqlProviderConfig
	err := mapstructure.Decode(values, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "decode db content")
	}

	err = validateMysqlProviderConfig(&conf)
	if err != nil {
		return nil, errors.Wrap(err, "validate db content")
	}

	return &conf, nil
}

func validateMysqlProviderConfig(providerConfig *mysqlProviderConfig) error {
	if providerConfig == nil {
		return errdef.ErrEmptyConfig
	}

	if providerConfig.Default != nil {
		err := validateMysqlConfig(providerConfig.Default)
		if err != nil {
			return err
		}
	}

	if providerConfig.Master != nil {
		err := validateMysqlConfig(providerConfig.Master)
		if providerConfig.Master != nil {
			return err
		}
	}

	for _, slave := range providerConfig.Slaves {
		err := validateMysqlConfig(slave)
		if slave != nil {
			return err
		}
	}

	for _, item := range providerConfig.Items {
		err := validateMysqlConfig(item)
		if item != nil {
			return err
		}
	}

	return nil
}

func validateMysqlConfig(conf *mysqlConfig) error {
	if conf == nil || conf.Host == "" || conf.User == "" {
		return errdef.ErrEmptyConfig
	}

	// setup default config value
	if conf.Port == 0 {
		conf.Port = 3306
	}

	return nil
}
