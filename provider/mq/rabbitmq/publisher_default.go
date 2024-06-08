package rabbitmq

import (
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
)

type rmqDefaultPublisherImpl struct {
	*rmqPublisherImpl
}

func newDefaultPublisher(config *RabbitMqConfig, logger intf.LoggerProvider) (intf.MessageQueuePublisher, error) {
	p, err := newPublisher(config, logger)
	if err != nil {
		return nil, err
	}

	return &rmqDefaultPublisherImpl{
		rmqPublisherImpl: p,
	}, nil
}

func (p *rmqDefaultPublisherImpl) Publish(topic string, messages [][]byte, delaySeconds ...int64) (err error) {
	t, err := newTopology(topic)
	if err != nil {
		return errors.Wrap(err, "new delay topology")
	}
	return p.publish(topic, messages, t)
}
