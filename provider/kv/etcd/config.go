package etcd

import (
	"github.com/hdget/hdsdk/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type ConfigEtcd struct {
	Url string `mapstructure:"url"`
}

// /////////////////////////////////////////////////////////////////
func parseConfig(rootConfiger types.Configer) (*ConfigEtcd, error) {
	if rootConfiger == nil {
		return nil, types.ErrEmptyConfig
	}

	data := rootConfiger.GetEtcdConfig()
	if data == nil {
		return nil, types.ErrEmptyConfig
	}

	values, ok := data.(map[string]interface{})
	if !ok {
		return nil, types.ErrInvalidConfig
	}

	var conf ConfigEtcd
	err := mapstructure.Decode(values, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "decode db content")
	}

	return &conf, nil
}

func validateConf(conf *ConfigEtcd) error {
	if conf == nil {
		return types.ErrEmptyConfig
	}

	if conf.Url == "" {
		return types.ErrInvalidConfig
	}

	return nil
}
