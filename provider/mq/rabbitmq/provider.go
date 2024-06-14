package rabbitmq

import (
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdsdk/v2/provider/mq"
)

// rabbitmqProvider
// Note: most codes comes from https://github.com/ThreeDotsLabs/watermill-amqp
type rabbitmqProvider struct {
	config *RabbitMqConfig
	logger intf.LoggerProvider
}

func New(configProvider intf.ConfigProvider, logger intf.LoggerProvider) (intf.MessageQueueProvider, error) {
	config, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	return &rabbitmqProvider{config: config, logger: logger}, nil
}

func (r rabbitmqProvider) Init(args ...any) error {
	//TODO implement me
	panic("implement me")
}

func (r rabbitmqProvider) Publisher(name string, args ...*mq.PublisherOption) (intf.MessageQueuePublisher, error) {
	option := mq.DefaultPublisherOption
	if len(args) > 0 {
		option = args[0]
	}

	publisherOptions := make([]publisherOption, 0)
	if option.PublishDelayMessage {
		publisherOptions = append(publisherOptions, withPublisherDelayTopology())
	}

	return newPublisher(name, r.config, r.logger, publisherOptions...)
}

func (r rabbitmqProvider) Subscriber(name string, args ...*mq.SubscriberOption) (intf.MessageQueueSubscriber, error) {
	option := mq.DefaultSubscriberOption
	if len(args) > 0 {
		option = args[0]
	}

	subscriberOptions := make([]subscriberOption, 0)
	if option.SubscribeDelayMessage {
		subscriberOptions = append(subscriberOptions, withSubscriberDelayTopology())
	}

	return newSubscriber(name, r.config, r.logger, subscriberOptions...)
}
