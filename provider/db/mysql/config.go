package mysql

import (
	"github.com/hdget/hdsdk/v1/errdef"
	"github.com/hdget/hdsdk/v1/intf"
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

	var providerConfig mysqlProviderConfig
	err := mapstructure.Decode(values, &providerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "decode mysql provider config")
	}

	err = providerConfig.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate mysql provider config")
	}

	return &providerConfig, nil
}

func (pc *mysqlProviderConfig) validate() error {
	if pc.Default != nil {
		err := pc.validateMysqlConfig(pc.Default)
		if err != nil {
			return err
		}
	}

	if pc.Master != nil {
		err := pc.validateMysqlConfig(pc.Master)
		if err != nil {
			return err
		}
	}

	for _, slave := range pc.Slaves {
		err := pc.validateMysqlConfig(slave)
		if err != nil {
			return err
		}
	}

	for _, item := range pc.Items {
		err := pc.validateMysqlExtraConfig(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pc *mysqlProviderConfig) validateMysqlConfig(conf *mysqlConfig) error {
	if conf == nil || conf.Host == "" || conf.User == "" {
		return errdef.ErrEmptyConfig
	}

	// setup default config value
	if conf.Port == 0 {
		conf.Port = 3306
	}

	return nil
}

func (pc *mysqlProviderConfig) validateMysqlExtraConfig(conf *mysqlConfig) error {
	if conf == nil || conf.Host == "" || pc.Master.User == "" || conf.Name == "" {
		return errdef.ErrEmptyConfig
	}

	// setup default config value
	if pc.Master.Port == 0 {
		pc.Master.Port = 3306
	}

	return nil
}
