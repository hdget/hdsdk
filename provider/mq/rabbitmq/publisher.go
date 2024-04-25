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
	config                  *RabbitMqConfig
	publishBindingsLock     sync.RWMutex
	publishBindingsPrepared map[string]struct{}
	closePublisher          func() error
	channelManager          channelManager
}

func NewPublisher(configProvider intf.ConfigProvider, logger intf.LoggerProvider) (intf.Publisher, error) {
	config, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	conn, err := newConnection(logger, config.Connection)
	if err != nil {
		return nil, fmt.Errorf("publisher create new connection: %w", err)
	}

	chanManager, err := newChannelManager(logger, conn, config.Connection.ChannelPoolSize)
	if err != nil {
		return nil, fmt.Errorf("create new channel pool: %w", err)
	}

	// Close the connection when the publisher is closed since this publisher owns the connection.
	closePublisher := func() error {
		logger.Debug("closing publisher connection")

		chanManager.Close()

		return conn.Close()
	}

	return &rmqPublisherImpl{
		connection:              conn,
		config:                  config,
		publishBindingsLock:     sync.RWMutex{},
		publishBindingsPrepared: make(map[string]struct{}),
		closePublisher:          closePublisher,
		channelManager:          chanManager,
	}, nil
}

// Publish publishes messages to AMQP broker.
// Publish is blocking until the broker has received and saved the message.
// Publish is always thread safe.
func (p *rmqPublisherImpl) Publish(topic string, messages ...[]byte) (err error) {
	exchangeKind := ExchangeKindDefault
	if p.config.UseExplicitExchange {
		exchangeKind = ExchangeKindExplicit
	}
	return p.publishByExchangeKind(exchangeKind, topic, 0, messages...)
}

// PublishDelay publishes delay messages to AMQP broker.
// Publish is blocking until the broker has received and saved the message.
// Publish is always thread safe.
func (p *rmqPublisherImpl) PublishDelay(topic string, delaySecond int64, messages ...[]byte) (err error) {
	return p.publishByExchangeKind(ExchangeKindDelay, topic, delaySecond, messages...)
}

// PublishDelay publishes delay messages to AMQP broker.
// Publish is blocking until the broker has received and saved the message.
// Publish is always thread safe.
func (p *rmqPublisherImpl) publishByExchangeKind(exchangeKind ExchangeKind, topic string, delaySecond int64, messages ...[]byte) error {
	if p.connection.IsClosed() {
		return errors.New("connection is closed while publish message")
	}

	if !p.connection.IsConnected() {
		return errors.New("connection is not established yet while publish message")
	}

	p.connectionWaitGroup.Add(1)
	defer p.connectionWaitGroup.Done()

	theChannel, err := p.channelManager.GetChannel()
	if err != nil {
		return errors.Wrap(err, "cannot open amqpChannel")
	}

	defer func() {
		if channelCloseErr := p.channelManager.CloseChannel(theChannel); channelCloseErr != nil {
			p.logger.Error("close AMQP channel", "err", channelCloseErr)
		}
	}()

	// get topology information
	topology := newTopology(topic, exchangeKind, p.config.ExchangeType)

	err = p.preparePublishBindings(topic, theChannel.AMQPChannel(), topology)
	if err != nil {
		return err
	}

	for _, msgPayload := range messages {
		if err = p.publishMessage(theChannel, topology, msgPayload, delaySecond); err != nil {
			return err
		}
	}

	return nil
}

func (p *rmqPublisherImpl) preparePublishBindings(topic string, amqpChannel *amqp.Channel, topology *Topology) error {
	p.publishBindingsLock.RLock()
	_, prepared := p.publishBindingsPrepared[topic]
	p.publishBindingsLock.RUnlock()

	if prepared {
		return nil
	}

	p.publishBindingsLock.Lock()
	defer p.publishBindingsLock.Unlock()

	// setup bindings
	if topology.exchangeName != "" {
		err := topology.declareExchange(amqpChannel)
		if err != nil {
			return errors.Wrap(err, "declare exchange while prepare publish bindings")
		}
	}

	p.publishBindingsPrepared[topic] = struct{}{}
	return nil
}

func (p *rmqPublisherImpl) publishMessage(channel channel, topology *Topology, msgPayload []byte, delaySecond int64) error {
	var headers amqp.Table
	if topology.exchangeKind == ExchangeKindDelay {
		if delaySecond <= 0 {
			return errors.New("delay must be greater than 0")
		}
		headers = map[string]interface{}{
			"x-delay": delaySecond * 1000, // message expire time in delay exchange, unit is mill seconds， provided by delay-message plugin
		}
	}

	err := channel.AMQPChannel().PublishWithContext(
		context.Background(),
		topology.exchangeName,
		topology.routingKey,
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
		p.logger.Trace("message published", "topology", topology, "msg", string(msgPayload))
		return nil
	}

	p.logger.Trace("message published, waiting for delivery confirmation", "topology", topology, "msg", string(msgPayload))
	if !channel.Delivered() {
		return fmt.Errorf("delivery not confirmed for message [%s]", string(msgPayload))
	}

	p.logger.Trace("delivery confirmed for message", "msg", string(msgPayload))
	return nil
}
