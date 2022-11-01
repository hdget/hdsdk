package neo4j

import (
	"github.com/hdget/hdsdk/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type ConfigNeo4j struct {
	VirtualUri  string             `mapstructure:"virtual_uri"`
	Username    string             `mapstructure:"username"`
	Password    string             `mapstructure:"password"`
	Servers     []*Neo4jServerConf `mapstructure:"servers"`
	MaxPoolSize int                `mapstructure:"max_pool_size"`
}

type Neo4jServerConf struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// /////////////////////////////////////////////////////////////////
func parseConfig(rootConfiger types.Configer) (*ConfigNeo4j, error) {
	if rootConfiger == nil {
		return nil, types.ErrEmptyConfig
	}

	data := rootConfiger.GetGraphConfig()
	if data == nil {
		return nil, types.ErrEmptyConfig
	}

	values, ok := data.(map[string]interface{})
	if !ok {
		return nil, types.ErrInvalidConfig
	}

	var conf ConfigNeo4j
	err := mapstructure.Decode(values, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "decode neo4j config")
	}

	return &conf, nil
}

func validateConf(providerType string, conf *ConfigNeo4j) error {
	if conf == nil {
		return types.ErrEmptyConfig
	}

	if conf.VirtualUri == "" || conf.Username == "" || conf.Password == "" {
		return types.ErrInvalidConfig
	}

	for _, server := range conf.Servers {
		if server.Host == "" || server.Port == 0 {
			return types.ErrInvalidConfig
		}
	}

	return nil
}
