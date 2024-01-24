package neo4j

import (
	"github.com/hdget/hdsdk/v1/errdef"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

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

func NewConfig(sdkConfiger intf.SdkConfiger) (*neo4jProviderConfig, error) {
	if sdkConfiger == nil {
		return nil, errdef.ErrEmptyConfig
	}

	values := sdkConfiger.GetNeo4jConfig()
	if len(values) == 0 {
		return nil, errdef.ErrEmptyConfig
	}

	var providerConfig neo4jProviderConfig
	err := mapstructure.Decode(values, &providerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "decode neo4j config")
	}

	err = providerConfig.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate neo4j config")
	}

	return &providerConfig, nil
}

func (pc *neo4jProviderConfig) validate() error {
	if pc.VirtualUri == "" || pc.Username == "" || pc.Password == "" {
		return errdef.ErrInvalidConfig
	}

	for _, server := range pc.Servers {
		if server.Host == "" || server.Port == 0 {
			return errdef.ErrInvalidConfig
		}
	}

	// setup default config items
	if pc.MaxPoolSize == 0 {
		pc.MaxPoolSize = defaultMaxPoolSize
	}

	return nil
}
