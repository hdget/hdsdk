package mysql

import (
	"github.com/hdget/hdsdk/v1/errdef"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/pkg/errors"
)

type providerConfig struct {
	mysql *mysqlProviderConfig `mapstructure:"mysql"` // 日志配置
}

type mysqlProviderConfig struct {
	Default *instanceConfig   `mapstructure:"default"`
	Master  *instanceConfig   `mapstructure:"master"`
	Slaves  []*instanceConfig `mapstructure:"slaves"`
	Items   []*instanceConfig `mapstructure:"items"`
}

type instanceConfig struct {
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Timeout  int    `mapstructure:"timeout"`
}

func NewConfig(configProvider intf.ConfigProvider) (*mysqlProviderConfig, error) {
	if configProvider == nil {
		return nil, errdef.ErrInvalidConfig
	}

	var c providerConfig
	err := configProvider.UnmarshalProviderConfig(&c)
	if err != nil {
		return nil, errdef.ErrInvalidConfig
	}

	err = c.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate mysql provider config")
	}

	return c.mysql, nil
}

func (c *providerConfig) validate() error {
	if c.mysql == nil {
		return errdef.ErrEmptyConfig
	}

	if c.mysql.Default != nil {
		err := c.validateInstance(c.mysql.Default)
		if err != nil {
			return err
		}
	}

	if c.mysql.Master != nil {
		err := c.validateInstance(c.mysql.Master)
		if err != nil {
			return err
		}
	}

	for _, slave := range c.mysql.Slaves {
		err := c.validateInstance(slave)
		if err != nil {
			return err
		}
	}

	for _, item := range c.mysql.Items {
		err := c.validateExtraInstance(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *providerConfig) validateInstance(ic *instanceConfig) error {
	if ic == nil || ic.Host == "" || ic.User == "" {
		return errdef.ErrEmptyConfig
	}

	// setup default config value
	if ic.Port == 0 {
		ic.Port = 3306
	}

	return nil
}

func (c *providerConfig) validateExtraInstance(conf *instanceConfig) error {
	if conf == nil || conf.Host == "" || conf.Name == "" {
		return errdef.ErrEmptyConfig
	}

	// setup default config value
	if conf.Port == 0 {
		conf.Port = 3306
	}
	return nil
}
