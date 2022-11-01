// Package kafka
package kafka

import (
	"hdsdk/types"
)

type Kafka struct {
	Logger types.LogProvider
	Config *MqConfig
}

var _ types.Mq = (*Kafka)(nil)

func NewMq(providerType string, config *MqConfig, logger types.LogProvider) (types.Mq, error) {
	err := validateMqConfig(providerType, config)
	if err != nil {
		return nil, err
	}
	return &Kafka{Logger: logger, Config: config}, nil
}

func (k *Kafka) GetDefaultOptions() map[types.MqOptionType]types.MqOptioner {
	return map[types.MqOptionType]types.MqOptioner{
		types.MqOptionPublish: defaultPublishOption,
		types.MqOptionConsume: defaultConsumeOption,
	}
}
