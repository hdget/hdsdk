package rabbitmq

import (
	"context"
	"fmt"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
)

type rmqPublisherImpl struct {
	*connection
	logger                  intf.LoggerProvider
	config                  *RabbitMqConfig
	publishBindingsLock     sync.RWMutex
	publishBindingsPrepared map[string]struct{}
	closePublisher          func() error
	channelManager          channelManager
}

func newPublisher(config *RabbitMqConfig, logger intf.LoggerProvider) (*rmqPublisherImpl, error) {
	conn, err := newConnection(logger, config)
	if err != nil {
		return nil, fmt.Errorf("publisher create new connection: %w", err)
	}

	channelManager, err := newChannelManager(logger, conn, config.ChannelPoolSize)
	if err != nil {
		return nil, fmt.Errorf("create new channel pool: %w", err)
	}

	// Close the connection when the publisher is closed since this publisher owns the connection.
	closePublisher := func() error {
		logger.Debug("closing publisher connection")

		channelManager.Close()

		return conn.Close()
	}

	return &rmqPublisherImpl{
		connection:              conn,
		logger:                  logger,
		config:                  config,
		publishBindingsLock:     sync.RWMutex{},
		publishBindingsPrepared: make(map[string]struct{}),
		closePublisher:          closePublisher,
		channelManager:          channelManager,
	}, nil
}

// Publish publishes messages to AMQP broker.
// Publish is blocking until the broker has received and saved the message.
// Publish is always thread safe.
func (p *rmqPublisherImpl) publish(topic string, messages [][]byte, t *Topology, args ...int64) (err error) {
	if p.connection.IsClosed() {
		return errors.New("connection is closed while publish message")
	}

	if !p.connection.IsConnected() {
		return errors.New("connection is not established yet while publish message")
	}

	p.connection.Begin()
	defer p.connection.End()

	theChannel, err := p.channelManager.GetChannel()
	if err != nil {
		return errors.Wrap(err, "cannot open amqpChannel")
	}

	defer func() {
		if channelCloseErr := p.channelManager.CloseChannel(theChannel); channelCloseErr != nil {
			p.logger.Error("close AMQP channel", "err", channelCloseErr)
		}
	}()

	err = p.preparePublishBindings(topic, theChannel.AMQPChannel(), t)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		if err = p.publishMessage(theChannel, t, msg, args...); err != nil {
			return err
		}
	}

	return nil
}

func (p *rmqPublisherImpl) preparePublishBindings(topic string, amqpChannel *amqp.Channel, t *Topology) error {
	p.publishBindingsLock.RLock()
	_, prepared := p.publishBindingsPrepared[topic]
	p.publishBindingsLock.RUnlock()

	if prepared {
		return nil
	}

	p.publishBindingsLock.Lock()
	defer p.publishBindingsLock.Unlock()

	if t.ExchangeName != "" {
		err := t.DeclareExchange(amqpChannel)
		if err != nil {
			return errors.Wrap(err, "declare exchange while prepare publish bindings")
		}
	}

	p.publishBindingsPrepared[topic] = struct{}{}
	return nil
}

func (p *rmqPublisherImpl) publishMessage(channel channel, t *Topology, msgPayload []byte, args ...int64) error {
	var headers amqp.Table
	if t.ExchangeKind == ExchangeKindDelay {
		if len(args) == 0 {
			return errors.New("no delay seconds specified")
		}
		headers = map[string]interface{}{
			"x-delay": args[0] * 1000, // message expire time in delay exchange, unit is mill seconds， provided by delay-message plugin
		}
	}

	err := channel.AMQPChannel().PublishWithContext(
		context.Background(),
		t.ExchangeName,
		t.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Body:         msgPayload,
			Headers:      headers,
			DeliveryMode: amqp.Persistent, // message always set to be persistent
		},
	)
	if err != nil {
		return errors.Wrap(err, "cannot publish msg")
	}

	if !channel.DeliveryConfirmationEnabled() {
		p.logger.Trace("message published", "topology", t, "msg", string(msgPayload))
		return nil
	}

	p.logger.Trace("message published, waiting for delivery confirmation", "topology", t, "msg", string(msgPayload))
	if !channel.Delivered() {
		return fmt.Errorf("delivery not confirmed for message [%s]", string(msgPayload))
	}

	p.logger.Trace("delivery confirmed for message", "msg", string(msgPayload))
	return nil
}
