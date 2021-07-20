package rabbitmq

import (
	"github.com/hdget/sdk/types"
	"github.com/pkg/errors"
)

type ProducerClient struct {
	*BaseClient
	Config *ProducerConfig
}

func (rmq *RabbitMq) newProducerClient(name string, options map[types.MqOptionType]types.MqOptioner) (*ProducerClient, error) {
	// 获取匹配的路由配置
	config, err := rmq.getProducerConfig(name)
	if config == nil {
		return nil, errors.Wrap(err, name)
	}
	return &ProducerClient{
		BaseClient: rmq.newBaseClient(name, options),
		Config:     config,
	}, nil
}
