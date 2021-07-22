package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/hdget/sdk/types"
)

// BaseClient 消息队列客户端维护connection和channel
type BaseClient struct {
	logger  types.LogProvider
	name    string // 名字
	options map[types.MqOptionType]types.MqOptioner
	ctx     context.Context
}

func (k *Kafka) newBaseClient(name string, options map[types.MqOptionType]types.MqOptioner) *BaseClient {
	// 将sarama日志输出
	sarama.Logger = k.Logger.GetStdLogger()

	return &BaseClient{
		logger:  k.Logger,
		name:    name,
		options: options,
		ctx:     context.Background(),
	}
}
