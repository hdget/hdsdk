package rabbitmq

import (
	"github.com/hdget/sdk/types"
	"github.com/pkg/errors"
)

type ConsumerClient struct {
	*BaseClient
	Config *ConsumerConfig
}

func (rmq *RabbitMq) NewConsumerClient(name string, options map[types.MqOptionType]types.MqOptioner) (*ConsumerClient, error) {
	// 获取匹配的路由配置
	config, err := rmq.getConsumerConfig(name)
	if err != nil {
		return nil, errors.Wrap(err, name)
	}

	return &ConsumerClient{
		BaseClient: rmq.newBaseClient(name, options),
		Config:     config,
	}, nil
}

// 声明和绑定queue
// @return error
func (cc *ConsumerClient) setupQueue() (string, error) {
	// 尝试声明队列, 检查指定的queue是否存在
	option := getQueueOption(cc.Options)
	q, err := cc.Channel.QueueDeclarePassive(
		cc.Config.QueueName,
		option.Durable,
		option.AutoDelete,
		option.Exclusive,
		option.NoWait,
		option.Args,
	)
	// 如果queue不存在，尝试声明这里注意如果声明出错，会关闭channel
	if err != nil {
		// 因为之前出错会关闭channel, 这里需要重连
		err := cc.connect()
		if err != nil {
			return "", err
		}

		q, err = cc.Channel.QueueDeclare(
			cc.Config.QueueName,
			option.Durable,
			option.AutoDelete,
			option.Exclusive,
			option.NoWait,
			option.Args,
		)
		// 注意queue声明失败后会关闭channel
		if err != nil {
			return "", err
		}
	}

	for _, key := range cc.Config.RoutingKeys {
		err = cc.Channel.QueueBind(q.Name, key, cc.Config.ExchangeName, true, nil)
		if err != nil {
			return "", errors.Wrap(err, "bind queue")
		}
	}

	return q.Name, nil
}
