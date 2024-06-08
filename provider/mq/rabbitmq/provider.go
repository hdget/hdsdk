package rabbitmq

import "github.com/hdget/hdsdk/v2/intf"

// rabbitmqProvider
// Note: most codes comes from https://github.com/ThreeDotsLabs/watermill-amqp
type rabbitmqProvider struct {
	config *RabbitMqConfig
	logger intf.LoggerProvider
}

var (
	_publisher       intf.MessageQueuePublisher
	_subscriber      intf.MessageQueueSubscriber
	_delayPublisher  intf.MessageQueuePublisher
	_delaySubscriber intf.MessageQueueSubscriber
)

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

func (r rabbitmqProvider) Publisher() (intf.MessageQueuePublisher, error) {
	var err error
	if _publisher == nil {
		_publisher, err = newDefaultPublisher(r.config, r.logger)
		if err != nil {
			return nil, err
		}
	}
	return _publisher, nil
}

func (r rabbitmqProvider) Subscriber() (intf.MessageQueueSubscriber, error) {
	var err error
	if _subscriber == nil {
		_subscriber, err = newDefaultSubscriber(r.config, r.logger)
		if err != nil {
			return nil, err
		}
	}
	return _subscriber, nil
}

func (r rabbitmqProvider) DelayPublisher(name string) (intf.MessageQueuePublisher, error) {
	var err error
	if _delayPublisher == nil {
		_delayPublisher, err = newDelayPublisher(r.config, r.logger, name)
		if err != nil {
			return nil, err
		}
	}

	return _delayPublisher, nil
}

func (r rabbitmqProvider) DelaySubscriber(name string) (intf.MessageQueueSubscriber, error) {
	var err error
	if _delaySubscriber == nil {
		_delaySubscriber, err = newDelaySubscriber(r.config, r.logger, name)
		if err != nil {
			return nil, err
		}
	}
	return _delaySubscriber, nil
}
