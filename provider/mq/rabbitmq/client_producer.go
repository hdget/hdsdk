package rabbitmq

import (
	"github.com/hdget/hdsdk/types"
)

type ProducerClient struct {
	*BaseClient
	//Config *ProducerConfig
}

func (rmq *RabbitMq) newProducerClient(options ...types.MqOptioner) (*ProducerClient, error) {
	return &ProducerClient{
		BaseClient: rmq.newBaseClient(options...),
	}, nil
}
