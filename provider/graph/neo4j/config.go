package neo4j

import (
	"github.com/hdget/hdsdk/v1/errdef"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/pkg/errors"
)

type providerConfig struct {
	neo4j *neo4jProviderConfig `mapstructure:"neo4j"` // 日志配置
}

type neo4jProviderConfig struct {
	VirtualUri  string             `mapstructure:"virtual_uri"`
	Username    string             `mapstructure:"username"`
	Password    string             `mapstructure:"password"`
	Servers     []*neo4jServerConf `mapstructure:"servers"`
	MaxPoolSize int                `mapstructure:"max_pool_size"`
}

type neo4jServerConf struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

const (
	defaultMaxPoolSize = 500
)

func NewConfig(configProvider intf.ConfigProvider) (*neo4jProviderConfig, error) {
	if configProvider == nil {
		return nil, errdef.ErrEmptyConfig
	}

	var c providerConfig
	err := configProvider.UnmarshalProviderConfig(&c)
	if err != nil {
		return nil, err
	}

	err = c.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate neo4j config")
	}

	return c.neo4j, nil
}

func (c *providerConfig) validate() error {
	if c.neo4j == nil {
		return errdef.ErrEmptyConfig
	}

	if c.neo4j.VirtualUri == "" || c.neo4j.Username == "" || c.neo4j.Password == "" {
		return errdef.ErrInvalidConfig
	}

	for _, server := range c.neo4j.Servers {
		if server.Host == "" || server.Port == 0 {
			return errdef.ErrInvalidConfig
		}
	}

	// setup default config items
	if c.neo4j.MaxPoolSize == 0 {
		c.neo4j.MaxPoolSize = defaultMaxPoolSize
	}

	return nil
}
