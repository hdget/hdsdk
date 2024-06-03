package sqlx_mysql

import (
	"github.com/hdget/hdsdk/v2/errdef"
	"github.com/hdget/hdsdk/v2/intf"
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

const (
	configSection = "sdk.mysql"
)

func newConfig(configProvider intf.ConfigProvider) (*mysqlProviderConfig, error) {
	if configProvider == nil {
		return nil, errdef.ErrInvalidConfig
	}

	var c *mysqlProviderConfig
	err := configProvider.Unmarshal(&c, configSection)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, errdef.ErrEmptyConfig
	}

	err = c.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate mysql provider config")
	}

	return c, nil
}

func (c *mysqlProviderConfig) validate() error {
	if c.Default != nil {
		err := c.validateInstance(c.Default)
		if err != nil {
			return err
		}
	}

	if c.Master != nil {
		err := c.validateInstance(c.Master)
		if err != nil {
			return err
		}
	}

	for _, slave := range c.Slaves {
		err := c.validateInstance(slave)
		if err != nil {
			return err
		}
	}

	for _, item := range c.Items {
		err := c.validateExtraInstance(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *mysqlProviderConfig) validateInstance(ic *mysqlConfig) error {
	if ic == nil || ic.Host == "" || ic.User == "" {
		return errdef.ErrEmptyConfig
	}

	// setup default config value
	if ic.Port == 0 {
		ic.Port = 3306
	}

	return nil
}

func (c *mysqlProviderConfig) validateExtraInstance(ic *mysqlConfig) error {
	if ic == nil || ic.Host == "" || ic.Name == "" {
		return errdef.ErrEmptyConfig
	}

	// setup default config value
	if ic.Port == 0 {
		ic.Port = 3306
	}
	return nil
}
