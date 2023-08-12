package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/hdget/hdsdk/types"
)

// BaseClient 消息队列客户端维护connection和channel
type BaseClient struct {
	logger  types.LogProvider
	options map[types.MqOptionType]types.MqOptioner
	ctx     context.Context
}

func (k *Kafka) newBaseClient(options ...types.MqOptioner) *BaseClient {
	// 设置传入option的值
	allOptions := k.GetDefaultOptions()
	for _, option := range options {
		allOptions[option.GetType()] = option
	}

	// 将sarama日志输出
	sarama.Logger = k.Logger.GetStdLogger()

	return &BaseClient{
		logger:  k.Logger,
		options: allOptions,
		ctx:     context.Background(),
	}
}
