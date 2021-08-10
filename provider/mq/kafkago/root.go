// Package kafkago
// @Title  log capability of zerolog
// @Description  zerolog implementation of log capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package kafkago

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/hdget/hdsdk/provider/mq"
	"github.com/hdget/hdsdk/types"
)

type KafkaProvider struct {
	mq.BaseMqProvider
}

var (
	_ types.Provider   = (*KafkaProvider)(nil)
	_ types.MqProvider = (*KafkaProvider)(nil)
)

// Init implements types.Provider interface, used to initialize the capability
// @author	Ryan Fan	(2021-06-09)
// @param	baseconf.Configer	root config interface to extract config info
// @return	error
func (kp *KafkaProvider) Init(rootConfiger types.Configer, logger types.LogProvider, args ...interface{}) error {
	// 获取日志配置信息
	config, err := parseConfig(rootConfiger)
	if err != nil {
		return err
	}

	if config.Default != nil {
		kp.Default, err = NewMq(types.PROVIDER_TYPE_DEFAULT, config.Default, logger)
		if err != nil {
			logger.Error("initialize kafkago", "type", types.PROVIDER_TYPE_DEFAULT, "brokers", config.Default.Brokers, "err", err)
		} else {
			logger.Debug("initialize kafkago", "type", types.PROVIDER_TYPE_DEFAULT, "brokers", config.Default.Brokers, "err", err)
		}
	}

	// 额外的mq
	kp.Items = make(map[string]types.Mq)
	for _, otherConf := range config.Items {
		instance, err := NewMq(types.PROVIDER_TYPE_OTHER, otherConf, logger)
		if err != nil {
			logger.Error("initialize mq", "type", otherConf.Name, "brokers", otherConf.Brokers, "err", err)
			continue
		}

		kp.Items[otherConf.Name] = instance
		logger.Debug("initialize mq", "type", otherConf.Name, "brokers", otherConf.Brokers, "err", err)
	}

	return nil
}
