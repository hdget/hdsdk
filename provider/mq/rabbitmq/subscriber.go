package rabbitmq

import (
	"context"
	"fmt"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdsdk/v2/provider/mq"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"sync/atomic"
	"time"
)

type rmpSubscriberImpl struct {
	*connection
	config              *RabbitMqConfig
	closedChan          chan struct{}
	closeSubscriber     func() error
	subscriberWaitGroup *sync.WaitGroup
}

func NewSubscriber(configProvider intf.ConfigProvider, logger intf.LoggerProvider) (intf.Subscriber, error) {
	config, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	conn, err := newConnection(logger, config.Connection)
	if err != nil {
		return nil, fmt.Errorf("subscriber create new connection: %w", err)
	}

	var closed uint32
	closedChan := make(chan struct{})
	var subscriberWaitGroup sync.WaitGroup

	// Close the subscriber AND the connection when the subscriber is closed,
	// since this subscriber owns the connection.
	closeSubscriber := func() error {
		if !atomic.CompareAndSwapUint32(&closed, 0, 1) {
			// Already closed.
			return nil
		}

		logger.Debug("closing subscriber.")

		close(closedChan)

		subscriberWaitGroup.Wait()

		logger.Debug("closing connection.")

		return conn.Close()
	}

	return &rmpSubscriberImpl{
		conn,
		config,
		closedChan,
		closeSubscriber,
		&subscriberWaitGroup,
	}, nil
}

// Subscribe consumes messages from AMQP broker.
func (s *rmpSubscriberImpl) Subscribe(ctx context.Context, topic string) (<-chan *mq.Message, error) {
	if s.connection.IsClosed() {
		return nil, errors.New("AMQP connection is closed")
	}

	if !s.connection.IsConnected() {
		return nil, errors.New("not connected to AMQP")
	}

	exchangeKind := ExchangeKindDefault
	if s.config.UseExplicitExchange {
		exchangeKind = ExchangeKindExplicit
	}

	topology := newTopology(topic, exchangeKind, s.config.ExchangeType)

	out := make(chan *mq.Message)
	if err := s.prepareConsume(topology); err != nil {
		return nil, errors.Wrap(err, "failed to prepare consume")
	}

	s.subscriberWaitGroup.Add(1)
	s.connectionWaitGroup.Add(1)

	go func(ctx context.Context) {
		defer func() {
			close(out)
			s.logger.Info("stopped consuming from AMQP channel")
			s.connectionWaitGroup.Done()
			s.subscriberWaitGroup.Done()
		}()

	ReconnectLoop:
		for {
			s.logger.Debug("waiting for connected or closing in reconnect loop")

			// to avoid race conditions with <-s.connected
			select {
			case <-s.connection.closing:
				s.logger.Debug("stopping reconnect loop (already closing)")
				break ReconnectLoop
			case <-s.closedChan:
				s.logger.Debug("stopping reconnect Loop (subscriber closing)")
				break ReconnectLoop
			default:
				// not closing yet
			}

			select {
			case <-s.connection.connected:
				s.logger.Debug("connection established in reconnect loop")
				// runSubscriber blocks until connection fails or Close() is called
				s.runSubscriber(ctx, out, topology)
			case <-s.connection.closing:
				s.logger.Debug("stopping reconnect loop (closing)")
				break ReconnectLoop
			case <-ctx.Done():
				s.logger.Debug("stopping reconnect loop (ctx done)")
				break ReconnectLoop
			}

			time.Sleep(time.Millisecond * 100)
		}
	}(ctx)

	return out, nil
}

func (s *rmpSubscriberImpl) prepareConsume(topology *Topology) error {
	amqpChannel, err := s.openSubscribeChannel()
	if err != nil {
		return err
	}

	defer func() {
		if channelCloseErr := amqpChannel.Close(); channelCloseErr != nil {
			s.logger.Error("close AMQP channel", "err", err)
		}
	}()

	err = s.prepareConsumeBindings(amqpChannel, topology)
	if err != nil {
		return err
	}
	return nil
}

func (s *rmpSubscriberImpl) prepareConsumeBindings(amqpChannel *amqp.Channel, topology *Topology) error {
	err := topology.declareQueue(amqpChannel)
	if err != nil {
		return errors.Wrap(err, "declare queue when prepare consume bindings")
	}

	if topology.exchangeName == "" {
		err = topology.declareExchange(amqpChannel)
		if err != nil {
			return errors.Wrap(err, "declare exchange when prepare consume bindings")
		}
	}

	err = topology.bindQueue(amqpChannel)
	if err != nil {
		return errors.Wrap(err, "bind queue when prepare consume bindings")
	}
	return nil
}

func (s *rmpSubscriberImpl) runSubscriber(ctx context.Context, out chan *mq.Message, topology *Topology) {
	amqpChannel, err := s.openSubscribeChannel()
	if err != nil {
		s.logger.Error("failed to open channel", "err", err)
		return
	}
	defer func() {
		if err = amqpChannel.Close(); err != nil {
			s.logger.Error("failed to close channel", "err", err)
		}
	}()

	notifyCloseChannel := amqpChannel.NotifyClose(make(chan *amqp.Error, 1))
	sub := subscription{
		out:                out,
		logger:             s.logger,
		notifyCloseChannel: notifyCloseChannel,
		channel:            amqpChannel,
		queueName:          topology.queueName,
		closing:            s.closing,
		closedChan:         s.closedChan,
		config:             s.config,
	}

	s.logger.Info("starting consuming from AMQP channel")

	sub.ProcessMessages(ctx)
}

func (s *rmpSubscriberImpl) openSubscribeChannel() (*amqp.Channel, error) {
	if !s.connection.IsConnected() {
		return nil, errors.New("not connected to AMQP")
	}

	amqpChannel, err := s.amqpConnection.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "cannot open channel")
	}

	err = amqpChannel.Qos(s.config.PrefetchCount, 0, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set channel Qos")
	}

	return amqpChannel, nil
}
