package rabbitmq

import "github.com/hdget/hdsdk/v2/intf"

type rabbitmqProvider struct {
	config *RabbitMqConfig
	logger intf.LoggerProvider
}

// New initialize zerolog instance
func New(configProvider intf.ConfigProvider, logger intf.LoggerProvider) (intf.MqProvider, error) {
	config, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	return &rabbitmqProvider{config: config, logger: logger}, nil
}

func (p rabbitmqProvider) Init(args ...any) error {
	panic("implement me")
}

func (p rabbitmqProvider) NewPublisher() (intf.Publisher, error) {
	return newPublisher(p.config, p.logger)
}

func (p rabbitmqProvider) NewSubscriber() (intf.Subscriber, error) {
	return newSubscriber(p.config, p.logger)
}
