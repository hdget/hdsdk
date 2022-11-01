// Package rabbitmq
//  1. exchangeType: direct
//     *  simple模式
//     一对一，一个发送一个接收
//     *  work模式
//     和simple模式一样，不同在于work模式可以有多个消费者，work模式起到一个负载均衡的作用,不同worker轮询获取一条消息并进行处理
//     *  Routing模式(publishKey: routingKey + queueBind: routingKey)
//     在订阅模式下，一个消息可以被多个消费者消费，如果我们想指定某个消息由哪些消费者消费，我们就要采用Routing模式,
//     routing模式最大的特点是可以从生产端通过routingKey来指定的消费端来消费消息
//  2. exchangeType: fanout
//     *  订阅模式(exchangeType: fanout, routingKey: empty)
//     simple模式和work模式他们有一个共同的特点就是一个消息只能被一个消费者消费，如果需要一个消息被多个消费者消费，就需要订阅模式。
//     订阅模式的特点是一个消息被投递到多个队列，一个消息能被多个消费者获取。
//     过程是由生产者将消息发送到exchange(交换机）里，然后exchange通过一系列的规则发送到队列上，然后由绑定对应的消费者进行消息。
//  3. exchangeType: topic
//     *  topic模式(publishKey: routingKey + queueBind: routingKey)
//     当基于routing模式，通过指定通配符的方式来指定我们的消费者来消费消息
package rabbitmq

import (
	"hdsdk/types"
)

type RabbitMq struct {
	Logger types.LogProvider
	Config *MqConfig
}

var _ types.Mq = (*RabbitMq)(nil)

func NewMq(providerType string, config *MqConfig, logger types.LogProvider) (types.Mq, error) {
	err := validateMqConfig(providerType, config)
	if err != nil {
		return nil, err
	}
	return &RabbitMq{Logger: logger, Config: config}, nil
}

func (rmq *RabbitMq) GetDefaultOptions() map[types.MqOptionType]types.MqOptioner {
	return map[types.MqOptionType]types.MqOptioner{
		types.MqOptionQueue:    defaultQueueOption,
		types.MqOptionExchange: defaultExchangeOption,
		types.MqOptionPublish:  defaultPublishOption,
		types.MqOptionConsume:  defaultConsumeOption,
		types.MqOptionQos:      defaultQosOption,
	}
}
