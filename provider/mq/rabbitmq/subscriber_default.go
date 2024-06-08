package rabbitmq

import (
	"context"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdsdk/v2/provider/mq"
	"github.com/pkg/errors"
)

type rmqDefaultSubscriberImpl struct {
	*rmpSubscriberImpl
}

func newDefaultSubscriber(config *RabbitMqConfig, logger intf.LoggerProvider) (intf.MessageQueueSubscriber, error) {
	p, err := newSubscriber(config, logger)
	if err != nil {
		return nil, err
	}

	return &rmqDefaultSubscriberImpl{
		rmpSubscriberImpl: p,
	}, nil
}

func (s *rmqDefaultSubscriberImpl) Subscribe(ctx context.Context, topic string) (<-chan *mq.Message, error) {
	t, err := newTopology(topic)
	if err != nil {
		return nil, errors.Wrap(err, "new topology")
	}

	return s.subscribe(ctx, t)
}
