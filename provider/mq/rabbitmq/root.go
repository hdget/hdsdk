// Package rabbitmq log capability of zerolog
// @Description  zerolog implementation of log capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package rabbitmq

import (
	"github.com/hdget/sdk/provider/mq"
	"github.com/hdget/sdk/types"
)

type RabbitmqProvider struct {
	mq.BaseMqProvider
}

var (
	_ types.Provider   = (*RabbitmqProvider)(nil)
	_ types.MqProvider = (*RabbitmqProvider)(nil)
)

// @desc	implements types.Provider interface, used to initialize the capability
// @author	Ryan Fan	(2021-06-09)
// @param	baseconf.Configer	root config interface to extract config info
// @return	error
func (rp *RabbitmqProvider) Init(rootConfiger types.Configer, logger types.LogProvider, args ...interface{}) error {
	// 获取日志配置信息
	config, err := parseConfig(rootConfiger)
	if err != nil {
		return err
	}

	rp.Default, err = NewMq(types.PROVIDER_TYPE_DEFAULT, config.Default, logger)
	if err != nil {
		logger.Error("initialize mq", "type", types.PROVIDER_TYPE_DEFAULT, "host", config.Default.Host, "err", err)
	} else {
		logger.Debug("initialize mq", "type", types.PROVIDER_TYPE_DEFAULT, "host", config.Default.Host, "err", err)
	}

	// 额外的mq
	rp.Items = make(map[string]types.Mq)
	for _, otherConf := range config.Items {
		instance, err := NewMq(types.PROVIDER_TYPE_OTHER, otherConf, logger)
		if err != nil {
			logger.Error("initialize mq", "type", otherConf.Name, "host", otherConf.Host, "err", err)
			continue
		}

		rp.Items[otherConf.Name] = instance
		logger.Debug("initialize mq", "type", otherConf.Name, "host", otherConf.Host, "err", err)
	}

	return nil
}
