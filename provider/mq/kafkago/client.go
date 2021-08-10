package kafkago

import (
	"context"
	"github.com/hdget/hdsdk/types"
)

// 消息队列客户端维护connection和channel
type BaseClient struct {
	Logger types.LogProvider

	Name    string // 客户端名字
	Options map[types.MqOptionType]types.MqOptioner
	ctx     context.Context
}

func (k *Kafka) newBaseClient(name string, options map[types.MqOptionType]types.MqOptioner) *BaseClient {
	// 连接URL
	return &BaseClient{
		Logger:  k.Logger,
		Name:    name,
		Options: options,
		ctx:     context.Background(),
	}
}
