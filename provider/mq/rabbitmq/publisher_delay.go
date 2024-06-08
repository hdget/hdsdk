package rabbitmq

import (
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
)

type rmqDelayPublisherImpl struct {
	*rmqPublisherImpl
	name string
}

func newDelayPublisher(config *RabbitMqConfig, logger intf.LoggerProvider, name string) (intf.MessageQueuePublisher, error) {
	p, err := newPublisher(config, logger)
	if err != nil {
		return nil, err
	}

	return &rmqDelayPublisherImpl{
		rmqPublisherImpl: p,
		name:             name,
	}, nil
}

func (p *rmqDelayPublisherImpl) Publish(topic string, messages [][]byte, delaySeconds ...int64) (err error) {
	if p.name == "" {
		return errors.New("publisher name must specified")
	}

	t, err := newDelayTopology(p.name, topic)
	if err != nil {
		return errors.Wrap(err, "new delay topology")
	}
	return p.publish(topic, messages, t, delaySeconds...)
}
