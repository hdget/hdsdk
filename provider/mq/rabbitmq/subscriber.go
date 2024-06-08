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
	logger              intf.LoggerProvider
	config              *RabbitMqConfig
	closedChan          chan struct{}
	closeSubscriber     func() error
	subscriberWaitGroup *sync.WaitGroup
}

func newSubscriber(config *RabbitMqConfig, logger intf.LoggerProvider) (*rmpSubscriberImpl, error) {
	conn, err := newConnection(logger, config)
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
		connection:          conn,
		logger:              logger,
		config:              config,
		closedChan:          closedChan,
		closeSubscriber:     closeSubscriber,
		subscriberWaitGroup: &subscriberWaitGroup,
	}, nil
}

// Subscribe consumes messages from AMQP broker.
func (s *rmpSubscriberImpl) subscribe(ctx context.Context, t *Topology) (<-chan *mq.Message, error) {
	if s.connection.IsClosed() {
		return nil, errors.New("AMQP connection is closed")
	}

	if !s.connection.IsConnected() {
		return nil, errors.New("not connected to AMQP")
	}

	out := make(chan *mq.Message)
	if err := s.prepareConsume(t); err != nil {
		return nil, errors.Wrap(err, "failed to prepare consume")
	}

	s.subscriberWaitGroup.Add(1)
	s.connection.Begin()

	go func(ctx context.Context) {
		defer func() {
			close(out)
			s.logger.Info("stopped consuming from AMQP channel")
			s.connection.End()
			s.subscriberWaitGroup.Done()
		}()

	ReconnectLoop:
		for {
			s.logger.Debug("waiting for connected or closing in reconnect loop")

			// to avoid race conditions with <-s.connected
			select {
			case <-s.connection.Closing():
				s.logger.Debug("stopping reconnect loop (already closing)")
				break ReconnectLoop
			case <-s.closedChan:
				s.logger.Debug("stopping reconnect Loop (subscriber closing)")
				break ReconnectLoop
			default:
				// not closing yet
			}

			select {
			case <-s.connection.Connected():
				s.logger.Debug("connection established in reconnect loop")
				// runSubscriber blocks until connection fails or Close() is called
				s.runSubscriber(ctx, out, t)
			case <-s.connection.Closing():
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

func (s *rmpSubscriberImpl) prepareConsume(t *Topology) error {
	amqpChannel, err := s.openSubscribeChannel()
	if err != nil {
		return err
	}

	defer func() {
		if channelCloseErr := amqpChannel.Close(); channelCloseErr != nil {
			s.logger.Error("close AMQP channel", "err", err)
		}
	}()

	err = s.prepareConsumeBindings(amqpChannel, t)
	if err != nil {
		return err
	}
	return nil
}

func (s *rmpSubscriberImpl) prepareConsumeBindings(amqpChannel *amqp.Channel, t *Topology) error {
	err := t.DeclareQueue(amqpChannel)
	if err != nil {
		return errors.Wrap(err, "declare queue when prepare consume bindings")
	}

	if t.ExchangeName != "" {
		err = t.DeclareExchange(amqpChannel)
		if err != nil {
			return errors.Wrap(err, "declare exchange when prepare consume bindings")
		}

		err = t.BindQueue(amqpChannel)
		if err != nil {
			return errors.Wrap(err, "bind queue when prepare consume bindings")
		}
	}

	return nil
}

func (s *rmpSubscriberImpl) runSubscriber(ctx context.Context, out chan *mq.Message, t *Topology) {
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
		queueName:          t.QueueName,
		closing:            s.connection.Closing(),
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

	amqpChannel, err := s.AmqpConnection().Channel()
	if err != nil {
		return nil, errors.Wrap(err, "cannot open channel")
	}

	err = amqpChannel.Qos(s.config.PrefetchCount, 0, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set channel Qos")
	}

	return amqpChannel, nil
}
