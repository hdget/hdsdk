package rabbitmq

import (
	"context"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdsdk/v2/provider/mq"
	"github.com/hdget/hdutils/text"
	"github.com/pkg/errors"
)

type rmqDelaySubscriberImpl struct {
	*rmpSubscriberImpl
	name string
}

func newDelaySubscriber(config *RabbitMqConfig, logger intf.LoggerProvider, name string) (intf.MessageQueueSubscriber, error) {
	p, err := newSubscriber(config, logger)
	if err != nil {
		return nil, err
	}

	if text.CleanString(name) == "" {
		return nil, errors.New("subscriber name must specified")
	}

	return &rmqDelaySubscriberImpl{
		rmpSubscriberImpl: p,
		name:              name,
	}, nil
}

func (s *rmqDelaySubscriberImpl) Subscribe(ctx context.Context, topic string) (<-chan *mq.Message, error) {
	t, err := newDelayTopology(s.name, topic)
	if err != nil {
		return nil, errors.Wrap(err, "new delay topology")
	}
	return s.subscribe(ctx, t)
}
