package rabbitmq

import (
	"github.com/hdget/hdsdk/types"
)

type ProducerClient struct {
	*BaseClient
	//Config *ProducerConfig
}

func (rmq *RabbitMq) newProducerClient(options map[types.MqOptionType]types.MqOptioner) (*ProducerClient, error) {
	//// 获取匹配的路由配置
	//config, err := rmq.getProducerConfig(name)
	//if config == nil {
	//	return nil, errors.Wrap(err, name)
	//}
	return &ProducerClient{
		BaseClient: rmq.newBaseClient(options),
	}, nil
}
