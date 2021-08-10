// Package gokit microservice implemented by gokit
// @Author  Ryan Fan 2021-07-27
// @Update  Ryan Fan 2021-07-27
package gokit

import (
	"github.com/hdget/hdsdk/provider/ms"
	"github.com/hdget/hdsdk/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type GokitProviderConfig struct {
	Default *MicroServiceConfig   `mapstructure:"default"`
	Items   []*MicroServiceConfig `mapstructure:"items"`
}

type GokitProvider struct {
	ms.BaseMsProvider
}

var (
	_ types.Provider   = (*GokitProvider)(nil)
	_ types.MsProvider = (*GokitProvider)(nil)
)

// Init	implements types.Provider interface, used to initialize the capability
// @author	Ryan Fan	(2021-06-09)
// @param	types.Configer	root config interface to extract config info
// @return	error
func (gp *GokitProvider) Init(rootConfiger types.Configer, logger types.LogProvider, args ...interface{}) error {
	// 获取日志配置信息
	config, err := parseProviderConfig(rootConfiger)
	if err != nil {
		return err
	}

	gp.Default, err = NewMicroService(logger, config.Default)
	if err != nil {
		logger.Error("initialize microservice", "type", "default", "err", err)
	} else {
		logger.Debug("initialize microservice", "type", "default")
	}

	// 额外的microservice
	gp.Items = make(map[string]types.MicroService)
	for _, otherConf := range config.Items {
		instance, err := NewMicroService(logger, otherConf)
		if err != nil {
			logger.Error("initialize microservice", "type", otherConf.Name, "err", err)
			continue
		}

		gp.Items[otherConf.Name] = instance
		logger.Debug("initialize microservice", "type", otherConf.Name)
	}

	return nil
}

// 解析MicroService配置
func parseProviderConfig(rootConfiger types.Configer) (*GokitProviderConfig, error) {
	data := rootConfiger.GetMicroServiceConfig()
	if data == nil {
		return nil, types.ErrEmptyConfig
	}

	values, ok := data.(map[string]interface{})
	if !ok {
		return nil, types.ErrInvalidConfig
	}

	var conf GokitProviderConfig
	err := mapstructure.Decode(values, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "decode microservice config")
	}

	return &conf, nil
}
