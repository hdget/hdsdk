package mysql

import (
	"github.com/hdget/hdsdk/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type ConfigMysql struct {
	Default *MySqlConf   `mapstructure:"default"`
	Master  *MySqlConf   `mapstructure:"master"`
	Slaves  []*MySqlConf `mapstructure:"slaves"`
	Items   []*MySqlConf `mapstructure:"items"`
}

type MySqlConf struct {
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Timeout  int    `mapstructure:"timeout"`
}

///////////////////////////////////////////////////////////////////
func parseConfig(rootConfiger types.Configer) (*ConfigMysql, error) {
	if rootConfiger == nil {
		return nil, types.ErrEmptyConfig
	}

	data := rootConfiger.GetMysqlConfig()
	if data == nil {
		return nil, types.ErrEmptyConfig
	}

	values, ok := data.(map[string]interface{})
	if !ok {
		return nil, types.ErrInvalidConfig
	}

	var conf ConfigMysql
	err := mapstructure.Decode(values, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "decode db content")
	}

	return &conf, nil
}

func validateConf(providerType string, conf *MySqlConf) error {
	if conf == nil {
		return types.ErrEmptyConfig
	}

	if conf.Host == "" || conf.Database == "" || conf.User == "" {
		return types.ErrInvalidConfig
	}

	if providerType == types.PROVIDER_TYPE_OTHER && conf.Name == "" {
		return types.ErrInvalidConfig
	}

	return nil
}
